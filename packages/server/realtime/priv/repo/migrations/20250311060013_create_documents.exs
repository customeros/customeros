defmodule Realtime.Repo.Migrations.CreateDocuments do
  use Ecto.Migration

  def change do
    create table(:documents, primary_key: false) do
      add(:id, :uuid, primary_key: true, default: fragment("gen_random_uuid()"))
      add(:name, :string)
      add(:body, :text)
      add(:lexical_state, :text)
      add(:tenant, :string)
      add(:user_id, :uuid)

      timestamps(type: :utc_datetime)
    end
  end
end
