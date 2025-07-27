package quiz

import (
	"context"
	"github.com/google/uuid"
	qpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (qs *QuizService) DeleteQuiz(ctx context.Context, req *qpb.DeleteQuizRequest) (*emptypb.Empty, error) {
	qs.logger.WithFields(logrus.Fields{
		"method": "DeleteQuiz",
		"req":    req,
	}).Info("Deleting quizzes")

	for _, id := range req.Ids {
		quizID, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quiz ID: %v", err)
		}

		// Check for existing questions
		questions, err := qs.store.GetQuestionsByQuizID(ctx, quizID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch quiz questions: %v", err)
		}

		if len(questions) > 0 && !req.ForceDelete {
			return nil, status.Errorf(codes.FailedPrecondition, "Quiz %s has associated questions. Use force_delete: true to force delete.", id)
		}

		if req.ForceDelete {
			// Delete answers
			for _, q := range questions {
				if err := qs.store.DeleteAnswersByQuestionId(ctx, uuid.MustParse(q.QuestionID.String())); err != nil {
					return nil, status.Errorf(codes.Internal, "failed to delete answers: %v", err)
				}
			}

			// Delete questions
			if err := qs.store.DeleteQuestionsByQuizId(ctx, quizID); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to delete questions: %v", err)
			}
		}

		// Delete quiz
		if err := qs.store.DeleteQuizById(ctx, quizID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to delete quiz: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}
