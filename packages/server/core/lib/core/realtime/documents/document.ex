defmodule Core.Realtime.Documents.Document do
  use Ecto.Schema
  import Ecto.Changeset

  @derive {Jason.Encoder,
           only: [
             :id,
             :name,
             :body,
             :lexical_state,
             :tenant,
             :user_id,
             :icon,
             :color,
             :organization_id,
             :inserted_at,
             :updated_at
           ]}
  @primary_key {:id, :binary_id, autogenerate: true}
  schema "documents" do
    field(:name, :string)
    field(:body, :string)
    field(:lexical_state, :string)
    field(:tenant, :string)
    field(:user_id, :binary_id)
    field(:icon, :string)
    field(:color, :string)
    field(:organization_id, :string, virtual: true)

    timestamps(type: :utc_datetime)
  end

  def changeset(document, attrs) do
    document
    |> cast(attrs, [:name, :body, :lexical_state, :tenant, :user_id, :icon, :color])
    |> validate_required([:name, :tenant, :body, :user_id, :icon, :color])
  end

  def update_changeset(document, attrs) do
    document
    |> cast(attrs, [:id, :name, :icon, :color])
    |> validate_required([:id, :name, :icon, :color])
  end
end
