import "@typespec/http";

using Http;

@service(#{ title: "Lesson Service" })
namespace LessonService {
  model Lesson {
    id: string; // UUID\n    title: string;\n    content: string | null;\n    materialType: string | null;\n    moduleId: string; // UUID\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  model UpdateLesson {}
  @route("/lessons")
  interface Lessons {
    @get list(): Lesson[];
    @get read(@path id: string): Lesson;
    @post create(@body body: Lesson): Lesson;
    @patch update(@path id: string, @body body: UpdateLesson): Lesson;
    @delete delete(@path id: string): void;
  }
}
