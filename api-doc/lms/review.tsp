import "@typespec/http";

using Http;

@service(#{ title: "Review Service" })
namespace ReviewService {
  model Review {
    id: string; // UUID\n    courseId: string; // UUID\n    userId: string; // UUID\n    title: string | null;\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/reviews")
  interface Reviews {
    @get list(): Review[];
    @get read(@path id: string): Review;
    @post create(@body body: Review): Review;
    @patch update(@path id: string, @body body: Partial<Review>): Review;
    @delete delete(@path id: string): void;
  }
}
