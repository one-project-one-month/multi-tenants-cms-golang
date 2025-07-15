import "@typespec/http";

using Http;

@service(#{ title: "Tenant Service" })
namespace TenantsService {
  model Tenant {
    id: string; // UUID\n    namespace: string;\n    cmsOwnerId: string; // UUID\n    createdAt: string; // timestamp\n    isActive: boolean;
  }

  @route("/tenants")
  interface Tenants {
    @get list(): Tenant[];
    @get read(@path id: string): Tenant;
    @post create(@body body: Tenant): Tenant;
    @patch update(@path id: string, @body body: Partial<Tenant>): Tenant;
    @delete delete(@path id: string): void;
  }
}
