defmodule RealtimeWeb.GenericMultiChannel do
  @moduledoc """
  This channel broadcasts sync events to multi-entity subscribers.
  It is used to dynamically handle multiple entity group types.
  """
  require Logger
  use RealtimeWeb, :channel

  @store Realtime.StoreManager

  def handle_join(
        entity_prefix,
        entity_id,
        %{"user_id" => user_id, "username" => username, "version" => client_version},
        socket
      ) do
    @store.prepare(entity_id)
    snapshot = @store.get_snapshot(entity_id)

    {_, latest_operations} =
      Enum.split(snapshot.operations, client_version)

    latest = %{
      version: snapshot.version,
      entity_id: snapshot.entity_id,
      operations: latest_operations
    }

    socket =
      socket
      |> assign(:user_id, user_id)
      |> assign(:username, username)
      |> assign(:entity_id, entity_id)
      |> assign(:entity_prefix, entity_prefix)

    send(self(), :after_join)
    {:ok, latest, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    {:noreply, socket}
  end

  @impl true
  def handle_in("sync_group_packet", payload, socket) do
    %{"payload" => %{"operation" => operation}} = payload

    broadcast!(socket, "sync_group_packet", operation)
    {:reply, {:ok, %{version: 0}}, socket}
  end

  @impl true
  def handle_in("sync_packet", payload, socket) do
    %{"payload" => %{"operation" => operation}} = payload

    entity_id = socket.assigns.entity_id

    sync_packet = %{
      version: 0,
      entity_id: entity_id,
      operation: operation
    }

    broadcast!(socket, "sync_packet", sync_packet)
    {:reply, {:ok, %{version: 0}}, socket}
  end

  @impl true
  def terminate(_, socket) do
    {:ok, socket}
  end

  @impl true
  def join(_, _, _) do
    {:ok, self()}
  end
end
