package quiz

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	qpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (qs *QuizService) CreateQuiz(ctx context.Context, req *qpb.CreateQuizRequest) (*qpb.QuizResponse, error) {
	qs.logger.WithFields(logrus.Fields{
		"method": "CreateQuiz",
		"params": req,
	}).Info("Creating quiz")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}
	organization := md.Get("x-organisation")
	if len(organization) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}
	organizationName := organization[0]

	// Validate ModuleID
	if req.ModuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "Module ID is required")
	}
	pgModuleID, err := global.ConvertStringToUUID(req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	fmt.Println("Org Name:", organizationName)
	fmt.Println("Module ID:", pgModuleID)
	// Check module exists
	args := repo.GetModuleByIDWithTenantParams{
		ModuleID:  uuid.MustParse(pgModuleID.String()),
		Namespace: organizationName,
	}
	module, err := qs.store.GetModuleByIDWithTenant(ctx, args)
	if err != nil {
		qs.logger.WithError(err).Error("Failed to fetch module")
		return nil, status.Errorf(codes.NotFound, "Module not found")
	}

	// Create Quiz
	quizResult, err := qs.store.CreateQuiz(ctx, db.CreateQuizParams{
		ModuleID: module.ModuleID,
		Title:    "Quiz for " + module.ModuleName,
	})
	if err != nil {
		qs.logger.WithError(err).Error("Failed to create quiz")
		return nil, status.Errorf(codes.Internal, "Unable to create quiz")
	}

	questionIDToQuestion := make(map[uuid.UUID]*qpb.Question)
	var allQuestions []*qpb.Question

	// Insert each question and its answers
	for _, q := range req.Quiz {
		questionResult, err := qs.store.CreateQuestion(ctx, db.CreateQuestionParams{
			Question: q.Question,
			QuizID:   quizResult.QuizID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to create question")
		}

		questionPB := &qpb.Question{
			QuestionId: questionResult.QuestionID.String(),
			QuizId:     quizResult.QuizID.String(),
			Question:   questionResult.Question,
			Answers:    []*qpb.Answer{},
		}
		allQuestions = append(allQuestions, questionPB)
		questionIDToQuestion[questionResult.QuestionID] = questionPB

		for _, a := range q.AnswerOptions {
			answerResult, err := qs.store.CreateAnswer(ctx, db.CreateAnswerParams{
				Answer:     a.Answer,
				IsCorrect:  pgtype.Bool{Bool: a.IsCorrect, Valid: true},
				QuestionID: questionResult.QuestionID,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to create answer")
			}

			questionPB.Answers = append(questionPB.Answers, &qpb.Answer{
				AnswerId:   answerResult.AnswerID.String(),
				QuestionId: answerResult.QuestionID.String(),
				Answer:     answerResult.Answer,
				IsCorrect:  answerResult.IsCorrect.Bool,
			})
		}
	}

	return &qpb.QuizResponse{
		Quiz: &qpb.Quiz{
			QuizId:    quizResult.QuizID.String(),
			ModuleId:  quizResult.ModuleID.String(),
			Title:     quizResult.Title,
			CreatedAt: timestamppb.New(quizResult.CreatedAt.Time),
			UpdatedAt: timestamppb.New(quizResult.UpdatedAt.Time),
		},
		Questions: allQuestions,
	}, nil
}
