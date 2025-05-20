defmodule Core.Realtime.YDoc do
  use Core.Realtime.YEcto, repo: Core.Repo, schema: Core.Realtime.Documents.DocumentWrite
end

defmodule Core.Realtime.EctoPersistence do
  @behaviour Yex.Sync.SharedDoc.PersistenceBehaviour
  @impl true
  def bind(_state, doc_name, doc) do
    ecto_doc = Core.Realtime.YDoc.get_y_doc(doc_name)

    {:ok, new_updates} = Yex.encode_state_as_update(doc)
    Core.Realtime.YDoc.insert_update(doc_name, new_updates)

    Yex.apply_update(doc, Yex.encode_state_as_update!(ecto_doc))
  end

  @impl true
  def unbind(_state, _doc_name, _doc) do
  end

  @impl true
  def update_v1(_state, update, doc_name, _doc) do
    Core.Realtime.YDoc.insert_update(doc_name, update)
    :ok
  end
end
