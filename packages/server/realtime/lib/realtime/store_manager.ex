defmodule Realtime.StoreManager do
  require Logger
  use GenServer

  @moduledoc """
  StoreManager is a GenServer that manages changesets and operations
  """

  @name __MODULE__

  @initial_state %{
    entity_id: nil,
    version: 0,
    operations: []
  }

  def start_link(_) do
    Logger.info("Store Manager started")
    GenServer.start_link(__MODULE__, [], name: @name)
  end

  def stop(pid) do
    Logger.info("Store Manager stopped")
    GenServer.stop(pid)
  end

  def prepare(entity_id), do: GenServer.call(@name, {:prepare, entity_id})
  def update(entity_id, change), do: GenServer.call(@name, {:update, entity_id, change})

  def set(entity_id, history, version),
    do: GenServer.call(@name, {:set_history, entity_id, history, version})

  def get_snapshot(entity_id), do: GenServer.call(@name, {:get_snapshot, entity_id})

  @impl true
  def init(_) do
    {:ok, @initial_state}
  end

  @impl true
  def handle_call({:prepare, entity_id}, _from, _state) do
    state = init_state(entity_id)
    {:reply, state, state}
  end

  @impl true
  def handle_call({:update, entity_id, change}, _from, state) do
    case get_state(entity_id) do
      nil ->
        Logger.error("State not found for entity_id: #{entity_id}")
        {:reply, {:error, "State not found"}, state}

      state ->
        new_state = %{
          entity_id: entity_id,
          version: state.version + 1,
          operations: state.operations ++ [change]
        }

        set_state(entity_id, new_state)
        {:reply, new_state, new_state}
    end
  end

  @impl true
  def handle_call({:set, entity_id, history, version}, _from, state) do
    case get_state(entity_id) do
      nil ->
        Logger.error("State not found for entity_id: #{entity_id}")
        {:reply, {:error, "State not found"}, state}

      _state ->
        new_state = %{
          entity_id: entity_id,
          version: version,
          operations: history
        }

        set_state(entity_id, new_state)
        {:reply, new_state, new_state}
    end
  end

  @impl true
  def handle_call({:get_snapshot, entity_id}, _from, state) do
    case get_state(entity_id) do
      nil ->
        Logger.error("State not found for entity_id: #{entity_id}")
        {:reply, {:error, "State not found"}, state}

      state ->
        # dbg(state)
        {:reply, state, state}
    end
  end

  # Private methods
  # ----------
  defp init_state(entity_id) do
    ensure_store_exists()

    case get_state(entity_id) do
      nil ->
        state = make_default_state(entity_id)
        :ets.insert(:obj_store, {entity_id, state})
        state

      state ->
        state
    end
  end

  defp ensure_store_exists() do
    case :ets.info(:obj_store) do
      :undefined -> :ets.new(:obj_store, [:named_table])
      _store -> :ok
    end
  end

  defp get_state(entity_id) do
    case :ets.lookup(:obj_store, entity_id) do
      [{^entity_id, state}] -> state
      _ -> nil
    end
  end

  defp set_state(entity_id, new_state) do
    case :ets.lookup(:obj_store, entity_id) do
      [{^entity_id, _prev_store}] -> :ets.insert(:obj_store, {entity_id, new_state})
      _ -> :ets.insert(:obj_store, {entity_id, new_state})
    end
  end

  defp make_default_state(entity_id) do
    %{
      entity_id: entity_id,
      version: 0,
      operations: []
    }
  end
end
