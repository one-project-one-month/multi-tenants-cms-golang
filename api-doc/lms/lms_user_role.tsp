import "@typespec/http";

using Http;

@service(#{ title: "UserRole Service" })
namespace LMSUserRoleService {
  model UserRole {
    id: string; // UUID\n    name: string; // lms_role_type
  }

  @route("/roles")
  interface UserRoles {
    @get list(): UserRole[];
    @get read(@path id: string): UserRole;
    @post create(@body body: UserRole): UserRole;
    @patch update(@path id: string, @body body: Partial<UserRole>): UserRole;
    @delete delete(@path id: string): void;
  }
}
