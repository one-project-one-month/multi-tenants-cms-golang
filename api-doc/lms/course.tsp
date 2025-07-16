import "@typespec/http";

using Http;

@service(#{ title: "Course Service" })
namespace CourseService {
  model Course {
    id: string; // UUID\n    title: string;\n    description: string | null;\n    instructorId: string; // UUID\n    overallRating: int32 | null;\n    categoryId: string; // UUID\n    status: string; // course_status\n    durationDayCount: int32 | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp\n    ownedBy: string; // UUID
  }

  @route("/courses")
  interface Courses {
    @get list(): Course[];
    @get read(@path id: string): Course;
    @post create(@body body: Course): Course;
    @patch update(@path id: string, @body body: Partial<Course>): Course;
    @delete delete(@path id: string): void;
  }
}
