import "@typespec/http";

using Http;

@service(#{ title: "Certificate Service" })
namespace CertificateService {
  model Certificate {
    id: string; // UUID\n    enrollmentId: string; // UUID\n    issueDate: string; // timestamp\n    certificateUrl: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/certificates")
  interface Certificates {
    @get list(): Certificate[];
    @get read(@path id: string): Certificate;
    @post create(@body body: Certificate): Certificate;
    @patch update(
      @path id: string,
      @body body: Partial<Certificate>,
    ): Certificate;
    @delete delete(@path id: string): void;
  }
}
