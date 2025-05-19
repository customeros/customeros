defmodule :"Elixir.Realtime.Repo.Migrations.Add-color-to-document" do
  use Ecto.Migration

  def change do
    alter table(:documents) do
      add :color, :string
    end
  end
end
