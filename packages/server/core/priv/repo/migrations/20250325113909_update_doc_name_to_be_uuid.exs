defmodule Realtime.Repo.Migrations.UpdateDocNameToBeUuid do
  use Ecto.Migration

  def change do
    alter table(:document_writes) do
      remove :docName
      add :docName, references(:documents, type: :uuid, on_delete: :delete_all), null: false
    end
  end
end
