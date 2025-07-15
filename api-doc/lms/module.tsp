import "@typespec/http";

using Http;

@service(#{ title: "Module Service" })
namespace ModuleService {
  model Module {
    id: string; // UUID\n    name: string;\n    courseId: string; // UUID\n    description: string | null;\n    createdAt: string; // timestamp\n    updatedAt: string | null; // timestamp
  }

  @route("/modules")
  interface Modules {
    @get list(): Module[];
    @get read(@path id: string): Module;
    @post create(@body body: Module): Module;
    @patch update(@path id: string, @body body: Partial<Module>): Module;
    @delete delete(@path id: string): void;
  }
}
