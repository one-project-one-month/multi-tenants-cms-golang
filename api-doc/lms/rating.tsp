import "@typespec/http";

using Http;

@service(#{ title: "Rating Service" })
namespace RatingService {
  model Rating {
    id: string; // UUID\n    userId: string; // UUID\n    courseId: string; // UUID\n    ratingCount: int32;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/ratings")
  interface Ratings {
    @get list(): Rating[];
    @get read(@path id: string): Rating;
    @post create(@body body: Rating): Rating;
    @patch update(@path id: string, @body body: Partial<Rating>): Rating;
    @delete delete(@path id: string): void;
  }
}
