defmodule Core.Realtime.Documents.OrganizationDocument do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key false
  schema "organization_documents" do
    field(:organization_id, :binary_id)

    belongs_to(:document, Core.Realtime.Documents.Document,
      type: :binary_id,
      foreign_key: :document_id
    )
  end

  def changeset(organization_document, attrs) do
    organization_document
    |> cast(attrs, [:organization_id, :document_id])
    |> validate_required([:organization_id, :document_id])
  end
end
