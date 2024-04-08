defmodule CustomerOsRealtimeWeb.TestChannel do
  @moduledoc """
  This is the Channel that tracks Test view.
  """
  alias Delta.Op
  require Logger
  use CustomerOsRealtimeWeb, :channel

  @delta CustomerOsRealtime.DeltaManager

  @impl true
  def join("org:" <> id, %{"user_id" => user_id, "username" => username}, socket) do
    @delta.prepare(id)

    socket =
      socket
      |> assign(:entity_id, id)
      |> assign(:user_id, user_id)
      |> assign(:username, username)

    send(self(), :after_join)
    {:ok, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    Logger.info("User #{socket.assigns.user_id} joined the channel")
    {:noreply, socket}
  end

  @impl true
  def handle_in("sync:" <> prop, payload, socket) do
    operations =
      Enum.map(payload, fn %{"type" => type, "payload" => payload} ->
        case type do
          "insert" -> Op.insert(payload)
          "delete" -> Op.delete(payload)
          "retain" -> Op.retain(payload)
          _ -> Logger.warning("Unknown action type: #{type}")
        end
      end)

    @delta.update(prop, operations)
    broadcast!(socket, "sync:" <> prop, @delta.get_contents(prop))
    {:noreply, socket}
  end

  @impl true
  def terminate(_, socket) do
    Logger.info("User #{socket.assigns.user_id} left the channel")
    {:ok, socket}
  end
end
