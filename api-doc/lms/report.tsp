import "@typespec/http";

using Http;

@service(#{ title: "Report Service" })
namespace ReportService {
  model Report {
    id: string; // UUID\n    reportName: string;\n    generatedByUserId: string; // UUID\n    generatedDate: string; // timestamp\n    dataSnapshot: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/reports")
  interface Reports {
    @get list(): Report[];
    @get read(@path id: string): Report;
    @post create(@body body: Report): Report;
    @patch update(@path id: string, @body body: Partial<Report>): Report;
    @delete delete(@path id: string): void;
  }
}
