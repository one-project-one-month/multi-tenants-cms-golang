import "@typespec/http";

using Http;

@service(#{ title: "Submission Service" })
namespace SubmissionService {
  model Submission {
    id: string; // UUID\n    assignmentId: string; // UUID\n    studentId: string; // UUID\n    submittedAt: string; // timestamp\n    fileUrl: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/submissions")
  interface Submissions {
    @get list(): Submission[];
    @get read(@path id: string): Submission;
    @post create(@body body: Submission): Submission;
    @patch update(
      @path id: string,
      @body body: Partial<Submission>,
    ): Submission;
    @delete delete(@path id: string): void;
  }
}
