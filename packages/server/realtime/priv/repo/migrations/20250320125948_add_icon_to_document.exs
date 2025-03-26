defmodule Realtime.Repo.Migrations.AddIconToDocument do
  use Ecto.Migration

  def change do
    alter table(:documents) do
      add :icon, :string
    end
  end
end
