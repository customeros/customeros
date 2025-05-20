defmodule Web.Graphql.Resolvers.DocumentResolver do
  alias Core.Realtime.Documents

  def list_documents(_parent, %{organization_id: id}, %{
        context: ctx
      }) do
    {:ok, Documents.list_by_organization(id, ctx.tenant)}
  end

  def get_document(_parent, %{id: id}, _ctx) do
    Documents.get_document(id)
  end

  def create_document(
        _parent,
        %{input: input},
        _ctx
      ) do
    {:ok, %{document: document}} = Documents.create_document(input)
    {:ok, document}
  end

  def update_document(_parent, args, _ctx) do
    Documents.update_document(args.input)
  end

  def delete_document(_parent, %{id: id}, _ctx) do
    Documents.delete_document(id)
  end
end
