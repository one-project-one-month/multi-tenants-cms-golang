import "@typespec/http";

using Http;

@service(#{ title: "NamespaceConsumer Service" })
namespace NamespaceConsumerService {
  model NamespaceConsumer {
    id: string; // UUID\n    userId: string; // UUID\n    namespace: string;\n    joinedDate: string; // timestamp\n    isActive: boolean;
  }

  @route("/namespace-consumers")
  interface NamespaceConsumers {
    @get list(): NamespaceConsumer[];
    @get read(@path id: string): NamespaceConsumer;
    @post create(@body body: NamespaceConsumer): NamespaceConsumer;
    @patch update(
      @path id: string,
      @body body: Partial<NamespaceConsumer>,
    ): NamespaceConsumer;
    @delete delete(@path id: string): void;
  }
}
