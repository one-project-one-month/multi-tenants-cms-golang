import "@typespec/http";

using Http;

@service(#{ title: "User Service" })
namespace LMSUserService {
  model User {
    id: string; // UUID\n    email: string;\n    password: string;\n    roleId: string; // UUID\n    tenantId: string | null;\n    address: string | null;\n    phoneNumber: string | null;\n    registrationDate: string | null; // date\n    createdAt: string; // timestamp\n    updatedAt: string; // timestamp
  }

  @route("/users")
  interface Users {
    @get list(): User[];
    @get read(@path id: string): User;
    @post create(@body body: User): User;
    @patch update(@path id: string, @body body: Partial<User>): User;
    @delete delete(@path id: string): void;
  }
}
