import "@typespec/http";

using Http;

@service(#{ title: "Enrollment Service" })
namespace EnrollmentService {
  model Enrollment {
    id: string; // UUID\n    studentId: string; // UUID\n    courseId: string; // UUID\n    enrollmentDate: string; // timestamp\n    progress: float32 | null;\n    status: string; // enrollment_type\n    dueDate: string | null; // timestamp\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/enrollments")
  interface Enrollments {
    @get list(): Enrollment[];
    @get read(@path id: string): Enrollment;
    @post create(@body body: Enrollment): Enrollment;
    @patch update(
      @path id: string,
      @body body: Partial<Enrollment>,
    ): Enrollment;
    @delete delete(@path id: string): void;
  }
}
