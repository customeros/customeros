defmodule Realtime.DocumentWrite do
  use Ecto.Schema
  import Ecto.Changeset

  schema "document_writes" do
    field(:value, :binary)
    field(:version, Ecto.Enum, values: [:v1, :v1_sv])
    field(:docName, :string)

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(yjs_doc, attrs) do
    yjs_doc
    |> cast(attrs, [:docName, :value, :version])
    |> validate_required([:docName, :value, :version])
  end
end
