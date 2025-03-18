defmodule RealtimeWeb.Graphql.Resolvers.DocumentResolver do
  alias Realtime.Documents

  def list_documents(_parent, %{organization_id: id}, %{context: %{tenant: tenant}}) do
    {:ok, Documents.list_by_organization(id, tenant)}
  end

  def create_document(
        _parent,
        %{
          input: %{
            name: _name,
            body: _body,
            user_id: _user_id,
            tenant: _tenant,
            lexical_state: _lexical_state,
            organization_id: _organization_id
          }
        } = args,
        _ctx
      ) do
    dbg(args.input)
    {:ok, %{document: document}} = Documents.create_document(args.input)
    {:ok, document}
  end
end
