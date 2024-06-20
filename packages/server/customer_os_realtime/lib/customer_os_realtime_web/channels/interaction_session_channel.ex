defmodule CustomerOsRealtimeWeb.InteractionSessionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSession entity subscribers.
  """
  require Logger
  use CustomerOsRealtimeWeb, :channel

  @store CustomerOsRealtime.StoreManager

  @impl true
  def join(
        "InteractionSession:" <> entity_id,
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

    send(self(), :after_join)
    {:ok, latest, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    Logger.info("User #{socket.assigns.user_id} joined the channel")
    {:noreply, socket}
  end

  @impl true
  def handle_in("sync_packet", payload, socket) do
    %{"payload" => %{"operation" => operation}} = payload
    entity_id = socket.assigns.entity_id

    # @store.update(entity_id, operation)
    snapshot = @store.get_snapshot(entity_id)

    sync_packet = %{
      version: snapshot.version,
      entity_id: snapshot.entity_id,
      operation: operation
    }

    broadcast!(socket, "sync_packet", sync_packet)
    {:reply, {:ok, %{version: snapshot.version}}, socket}
  end

  @impl true
  def terminate(_, socket) do
    Logger.info("User #{socket.assigns.user_id} left the channel")
    {:ok, socket}
  end
end
