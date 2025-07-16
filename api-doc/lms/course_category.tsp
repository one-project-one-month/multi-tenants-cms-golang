import "@typespec/http";

using Http;

@service(#{ title: "CourseCategory Service" })
namespace CourseCategoryService {
  model CourseCategory {
    id: string; // UUID\n    name: string;\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/course-categories")
  interface CourseCategorys {
    @get list(): CourseCategory[];
    @get read(@path id: string): CourseCategory;
    @post create(@body body: CourseCategory): CourseCategory;
    @patch update(
      @path id: string,
      @body body: Partial<CourseCategory>,
    ): CourseCategory;
    @delete delete(@path id: string): void;
  }
}
