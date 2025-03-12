defmodule Realtime.Repo.Migrations.CreateYjsDocumentWrites do
  use Ecto.Migration

  def change do
    create table(:document_writes) do
      add(:docName, :string)
      add(:value, :binary)
      add(:version, :string)

      timestamps(type: :utc_datetime)
    end

    create(index("document_writes", [:docName, :version]))
  end
end
