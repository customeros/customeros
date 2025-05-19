defmodule Realtime.Repo.Migrations.CreateOrganizationDocuments do
  use Ecto.Migration

  def change do
    create table(:organization_documents, primary_key: false) do
      add :organization_id, :binary_id, null: false
      add :document_id, references(:documents, type: :uuid, on_delete: :delete_all), null: false
    end

    create index(:organization_documents, [:organization_id])
    create index(:organization_documents, [:organization_id, :document_id], unique: true)
  end
end
