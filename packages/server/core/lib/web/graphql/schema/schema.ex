defmodule Web.Graphql.Schema do
  use Absinthe.Schema

  import_types(Web.Graphql.DocumentTypes)

  query do
    import_fields(:document_queries)
  end

  mutation do
    import_fields(:document_mutations)
  end
end
