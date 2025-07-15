import "@typespec/http";

using Http;

@service(#{ title: "StudentQuiz Service" })
namespace StudentQuizService {
  model StudentQuiz {
    id: string; // UUID\n    studentId: string; // UUID\n    quizId: string; // UUID\n    score: int32 | null;\n    attempt: int32;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/student-quizzes")
  interface StudentQuizs {
    @get list(): StudentQuiz[];
    @get read(@path id: string): StudentQuiz;
    @post create(@body body: StudentQuiz): StudentQuiz;
    @patch update(
      @path id: string,
      @body body: Partial<StudentQuiz>,
    ): StudentQuiz;
    @delete delete(@path id: string): void;
  }
}
