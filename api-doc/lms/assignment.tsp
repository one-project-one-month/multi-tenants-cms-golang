import "@typespec/http";

using Http;

@service(#{ title: "Assignment Service" })
namespace AssignmentService {
  model Assignment {
    id: string; // UUID\n    courseId: string; // UUID\n    title: string;\n    instructions: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/assignments")
  interface Assignments {
    @get list(): Assignment[];
    @get read(@path id: string): Assignment;
    @post create(@body body: Assignment): Assignment;
    @patch update(
      @path id: string,
      @body body: Partial<Assignment>,
    ): Assignment;
    @delete delete(@path id: string): void;
  }
}
