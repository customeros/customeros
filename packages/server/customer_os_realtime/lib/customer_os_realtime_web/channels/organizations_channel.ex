defmodule CustomerOsRealtimeWeb.OrganizationsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organizations entity subscribers.
  """
  require Logger
  use CustomerOsRealtimeWeb, :channel

  @store CustomerOsRealtime.StoreManager

  @impl true
  def join(
        "Organizations:" <> entity_id,
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
  def join("Organizations:" <> _entity_id, _params, socket) do
    Logger.info("User #{socket.assigns.user_id} joining the Organizations channel")
    send(self(), :after_join)
    {:ok, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    Logger.info("User #{socket.assigns.user_id} joined the Organizations channel")
    {:noreply, socket}
  end

  def handle_info(%{event: "sync_group_packet", payload: payload}, socket) do
    dbg(payload)
    push(socket, "sync_group_packet", payload)
    {:noreply, socket}
  end

  @impl true
  def handle_in("sync_group_packet", payload, socket) do
    dbg(payload)
    %{"payload" => %{"operation" => operation}} = payload

    broadcast!(socket, "sync_group_packet", operation)
    {:reply, {:ok, %{version: 0}}, socket}
  end

  @impl true
  def handle_in("sync_packet", payload, socket) do
    %{"payload" => %{"operation" => operation}} = payload

    broadcast!(socket, "sync_packet", operation)
    {:reply, {:ok, %{version: 0}}, socket}
  end

  @impl true
  def terminate(_, socket) do
    Logger.info("User #{socket.assigns.user_id} left the Organizations channel")
    {:ok, socket}
  end
end
