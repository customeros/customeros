defmodule Web.Graphql.DocumentTypes do
  use Absinthe.Schema.Notation
  alias Web.Graphql.Resolvers.DocumentResolver

  object :document do
    field(:id, non_null(:id))
    field(:name, non_null(:string))
    field(:body, :string)
    field(:lexical_state, :string)
    field(:tenant, :string)
    field(:user_id, non_null(:id))
    field(:icon, :string)
    field(:color, :string)
    field(:organization_id, :string)
    field(:inserted_at, non_null(:string))
    field(:updated_at, non_null(:string))
  end

  input_object :create_document_input do
    field(:name, non_null(:string))
    field(:body, non_null(:string))
    field(:tenant, non_null(:string))
    field(:user_id, non_null(:id))
    field(:icon, non_null(:string))
    field(:color, non_null(:string))
    field(:lexical_state, :string)
    field(:organization_id, non_null(:id))
  end

  input_object :update_document_input do
    field(:id, non_null(:id))
    field(:name, non_null(:string))
    field(:icon, non_null(:string))
    field(:color, non_null(:string))
  end

  object :document_queries do
    field :organization_documents, list_of(:document) do
      arg(:organization_id, non_null(:id))
      resolve(&DocumentResolver.list_documents/3)
    end

    field :document, :document do
      arg(:id, non_null(:id))
      resolve(&DocumentResolver.get_document/3)
    end
  end

  object :document_mutations do
    field :create_document, :document do
      arg(:input, non_null(:create_document_input))
      resolve(&DocumentResolver.create_document/3)
    end

    field :update_document, :document do
      arg(:input, non_null(:update_document_input))
      resolve(&DocumentResolver.update_document/3)
    end

    field :delete_document, :document do
      arg(:id, non_null(:id))
      resolve(&DocumentResolver.delete_document/3)
    end
  end
end
