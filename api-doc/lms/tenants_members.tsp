import "@typespec/http";

using Http;

@service(#{ title: "TenantMember Service" })
namespace TenantsMembersService {
  model TenantMember {
    id: string; // UUID\n    userId: string; // UUID\n    tenantId: string; // UUID\n    joinedDate: string; // timestamp\n    isActive: boolean;
  }

  @route("/tenant-members")
  interface TenantMembers {
    @get list(): TenantMember[];
    @get read(@path id: string): TenantMember;
    @post create(@body body: TenantMember): TenantMember;
    @patch update(
      @path id: string,
      @body body: Partial<TenantMember>,
    ): TenantMember;
    @delete delete(@path id: string): void;
  }
}
