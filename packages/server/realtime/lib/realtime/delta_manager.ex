defmodule Realtime.DeltaManager do
  require Logger
  require Delta
  use GenServer

  @moduledoc """
  DeltaManager is a GenServer that manages changesets and operations
  on a document. It keeps track of the changes made to the document
  and allows for viewing history, undoing changes, and more.
  """
  @name __MODULE__

  @initial_state %{
    entity_id: nil,
    # Number of changes made to the document so far
    version: 0,

    # An up-to-date Delta with all changes applied, representing
    # the current state of the document
    contents: [],

    # The `inverted` versions of all changes performed on the
    # document (useful for viewing history or undo the changes)
    inverted_changes: []
  }

  # Public API
  # ----------

  def start_link(_) do
    Logger.info("Delta Manager started")
    GenServer.start_link(__MODULE__, [], name: @name)
  end

  def stop(pid) do
    Logger.info("Delta Manager stopped")
    GenServer.stop(pid)
  end

  def prepare(entity_id), do: GenServer.call(@name, {:prepare, entity_id})
  def update(change), do: GenServer.call(@name, {:update, change})
  def update(prop, change), do: GenServer.call(@name, {:update, prop, change})
  def get_contents(prop), do: GenServer.call(@name, {:get_contents, prop})
  def get_history(), do: GenServer.call(@name, :get_history)
  def undo(), do: GenServer.call(@name, :undo)

  # GenServer Callbacks
  # -------------------

  # Initialize the document with the default state
  @impl true
  def init(_) do
    {:ok, @initial_state}
  end

  @impl true
  def handle_call({:prepare, entity_id}, _from, _state) do
    state = init_state(entity_id)
    {:reply, state, state}
  end

  # Apply a given change to the document, updating its contents
  # and incrementing the version
  #
  # We also keep track of the inverted version of the change
  # which is useful for performing undo or viewing history
  @impl true
  def handle_call({:update, change}, _from, state) do
    inverted = Delta.invert(change, state.contents)

    state = %{
      version: state.version + 1,
      contents: Delta.compose(state.contents, change),
      inverted_changes: [inverted | state.inverted_changes]
    }

    {:reply, state.contents, state}
  end

  @impl true
  def handle_call({:update, prop, change}, _from, state) do
    entity_id = state.entity_id

    case get_state(entity_id) do
      # nil ->
      #   {:reply, {:error, "State: unknown error"}, prev_state}

      # :not_found ->
      #   {:reply, {:error, "State not found"}, prev_state}

      state ->
        prop_atom = String.to_atom(prop)

        case Map.get(state, prop_atom) do
          nil ->
            {:reply, {:error, "Property not found"}, state}

          entry ->
            inverted = Delta.invert(change, entry)

            next_entry = %{
              version: entry.version + 1,
              contents: Delta.compose(entry.contents, change),
              inverted_changes: [inverted | entry.inverted_changes]
            }

            next_state = Map.put(state, prop_atom, next_entry)

            set_state(entity_id, next_state)
            {:reply, next_entry.contents, next_state}
        end
    end

    # inverted = Delta.invert(change, state.contents)

    # state = %{
    #   version: state.version + 1,
    #   contents: Delta.compose(state.contents, change),
    #   inverted_changes: [inverted | state.inverted_changes]
    # }

    # entry = Map.get(state, String.to_atom(prop))
    # dbg(entry)

    # {:reply, Map.get(entry, :contents), state}
  end

  # Fetch the current contents of the document
  @impl true
  def handle_call({:get_contents, prop}, _from, state) do
    case get_state(state.entity_id) do
      nil ->
        {:reply, {:error, "State: unknown error"}, state}

      :not_found ->
        {:reply, {:error, "State not found"}, state}

      state ->
        prop_atom = String.to_atom(prop)

        case Map.get(state, prop_atom) do
          nil ->
            {:reply, {:error, "Property not found"}, state}

          entry ->
            {:reply, entry.contents, state}
        end
    end
  end

  # Revert the applied changes one by one to see how the
  # document transformed over time
  @impl true
  def handle_call(:get_history, _from, state) do
    current = {state.version, state.contents}

    history =
      Enum.scan(state.inverted_changes, current, fn inverted, {version, contents} ->
        contents = Delta.compose(contents, inverted)
        {version - 1, contents}
      end)

    {:reply, [current | history], state}
  end

  # Don't undo when document is already empty
  @impl true
  def handle_call(:undo, _from, %{version: 0} = state) do
    {:reply, state.contents, state}
  end

  # Revert the last change, removing it from our stack and
  # updating the contents
  @impl true
  def handle_call(:undo, _from, state) do
    [last_change | changes] = state.inverted_changes

    state = %{
      version: state.version - 1,
      contents: Delta.compose(state.contents, last_change),
      inverted_changes: changes
    }

    {:reply, state.contents, state}
  end

  # Private methods
  # ----------
  defp init_state(entity_id) do
    ensure_store_exists()

    case get_state(entity_id) do
      nil ->
        state = make_default_state(entity_id)
        :ets.insert(:store, {entity_id, state})
        state

      state ->
        state
    end
  end

  defp ensure_store_exists() do
    case :ets.info(:store) do
      :undefined -> :ets.new(:store, [:named_table])
      _store -> :ok
    end
  end

  defp get_state(entity_id) do
    case :ets.lookup(:store, entity_id) do
      [{^entity_id, state}] -> state
      _ -> nil
    end
  end

  defp set_state(entity_id, new_state) do
    case :ets.lookup(:store, entity_id) do
      [{^entity_id, _prev_store}] -> :ets.insert(:store, {entity_id, new_state})
      _ -> :ets.insert(:store, {entity_id, new_state})
    end
  end

  defp make_default_state(entity_id) do
    %{
      entity_id: entity_id,
      name: %{
        version: 0,
        contents: [],
        inverted_changes: []
      },
      description: %{
        version: 0,
        contents: [],
        inverted_changes: []
      }
    }
  end
end
