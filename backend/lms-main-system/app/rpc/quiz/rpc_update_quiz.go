package quiz

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	qpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (qs *QuizService) UpdateQuiz(ctx context.Context, req *qpb.UpdateQuizRequest) (*qpb.Question, error) {
	qs.logger.WithFields(logrus.Fields{
		"method": "UpdateQuiz",
		"params": req,
	}).Info("Updating quiz")

	// Validate input
	if req.QuestionId == "" {
		return nil, status.Error(codes.InvalidArgument, "Question Id is required")
	}
	if req.Question == "" {
		return nil, status.Error(codes.InvalidArgument, "Question is required")
	}

	questionID, err := global.ConvertStringToUUID(req.QuestionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid question ID")
	}
	qUUID := uuid.MustParse(questionID.String())
	// Update the question
	updateQ := repo.UpdateQuestionParams{
		Question:   req.Question,
		QuestionID: qUUID,
	}

	_, err = qs.store.UpdateQuestion(ctx, updateQ)
	if err != nil {
		qs.logger.WithError(err).Error("Failed to update question")
		return nil, status.Error(codes.Internal, "Failed to update question")
	}

	// Update each answer option
	for _, a := range req.AnswerOptions {
		ua, err := qs.store.UpdateAnswer(ctx, repo.UpdateAnswerParams{
			Answer:    a.Answer,
			IsCorrect: pgtype.Bool{Bool: a.IsCorrect, Valid: true},
			AnswerID:  global.ConvertStringToGoogleUUID(a.AnswerId),
		})
		if err != nil {
			qs.logger.WithError(err).Errorf("Failed to update answer %s", ua.AnswerID)
			return nil, status.Errorf(codes.Internal, "Failed to update answer: %s", ua.AnswerID)
		}
	}

	// Fetch updated data to return
	question, err := qs.store.GetQuestionByID(ctx, qUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve updated question")
	}

	answers, err := qs.store.GetAnswersByQuestionID(ctx, qUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve updated answers")
	}

	responseAnswers := make([]*qpb.Answer, len(answers))
	for _, a := range answers {
		responseAnswers = append(responseAnswers, &qpb.Answer{
			AnswerId:   a.AnswerID.String(),
			Answer:     a.Answer,
			QuestionId: a.QuestionID.String(),
			IsCorrect:  a.IsCorrect.Bool,
			CreatedAt:  timestamppb.New(a.CreatedAt.Time),
			UpdatedAt:  timestamppb.New(a.CreatedAt.Time),
		})
	}

	return &qpb.Question{
		QuestionId: question.QuestionID.String(),
		QuizId:     question.QuizID.String(),
		Question:   question.Question,
		Answers:    responseAnswers,
	}, nil
}
