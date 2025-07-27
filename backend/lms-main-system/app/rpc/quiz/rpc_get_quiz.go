package quiz

import (
	"context"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	qpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (qs *QuizService) GetQuizzesByModule(ctx context.Context, req *qpb.GetQuizByModuleRequest) (*qpb.GetQuizByModuleResponse, error) {
	qs.logger.WithFields(logrus.Fields{
		"method": "GetQuizzesByModule",
		"params": req,
	}).Info("Fetching quizzes by module")

	// Validate input
	if req.ModuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "Module ID is required")
	}
	pgModuleID, err := global.ConvertStringToUUID(req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Module ID")
	}
	// Get quizzes by module ID
	quizzes, err := qs.store.GetQuizzesByModule(ctx, uuid.MustParse(pgModuleID.String()))
	if err != nil {
		qs.logger.WithError(err).Error("Failed to fetch quizzes")
		return nil, status.Errorf(codes.Internal, "Failed to fetch quizzes")
	}
	var quizResponses []*qpb.QuizResponse

	for _, quiz := range quizzes {
		questions, err := qs.store.GetQuestionsByQuizID(ctx, quiz.QuizID)
		if err != nil {
			qs.logger.WithError(err).Error("Failed to fetch questions")
			return nil, status.Errorf(codes.Internal, "Failed to fetch questions for quiz %s", quiz.QuizID)
		}

		var questionResponses []*qpb.Question

		for _, question := range questions {
			answers, err := qs.store.GetAnswersByQuestionID(ctx, question.QuestionID)
			if err != nil {
				qs.logger.WithError(err).Error("Failed to fetch answers")
				return nil, status.Errorf(codes.Internal, "Failed to fetch answers for question %s", question.QuestionID)
			}

			var answerResponses []*qpb.Answer
			for _, ans := range answers {
				answerResponses = append(answerResponses, &qpb.Answer{
					AnswerId:   ans.AnswerID.String(),
					QuestionId: ans.QuestionID.String(),
					Answer:     ans.Answer,
					IsCorrect:  ans.IsCorrect.Bool,
				})
			}

			questionResponses = append(questionResponses, &qpb.Question{
				QuestionId: question.QuestionID.String(),
				QuizId:     question.QuizID.String(),
				Question:   question.Question,
				Answers:    answerResponses,
			})
		}

		quizResponses = append(quizResponses, &qpb.QuizResponse{
			Quiz: &qpb.Quiz{
				QuizId:    quiz.QuizID.String(),
				ModuleId:  quiz.ModuleID.String(),
				Title:     quiz.Title,
				CreatedAt: timestamppb.New(quiz.CreatedAt.Time),
				UpdatedAt: timestamppb.New(quiz.UpdatedAt.Time),
			},
			Questions: questionResponses,
		})
	}

	return &qpb.GetQuizByModuleResponse{
		Quizzes: quizResponses,
	}, nil
}
