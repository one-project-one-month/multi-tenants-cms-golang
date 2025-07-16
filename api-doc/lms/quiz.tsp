import "@typespec/http";

using Http;

@service(#{ title: "Quiz Service" })
namespace QuizService {
  model Quiz {
    id: string; // UUID\n    question: string;\n    answer: string;\n    moduleId: string; // UUID\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/quizzes")
  interface Quizs {
    @get list(): Quiz[];
    @get read(@path id: string): Quiz;
    @post create(@body body: Quiz): Quiz;
    @patch update(@path id: string, @body body: Partial<Quiz>): Quiz;
    @delete delete(@path id: string): void;
  }
}
