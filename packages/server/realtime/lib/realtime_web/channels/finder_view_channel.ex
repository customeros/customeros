defmodule RealtimeWeb.FinderChannel do
  @moduledoc """
  This is the Channel that tracks Finder view.
  """
  require Logger
  use RealtimeWeb, :channel
  alias Realtime.ColorManager
  alias RealtimeWeb.Presence

  @impl true
  def join("finder:" <> _organization_id, %{"user_id" => user_id, "username" => username}, socket) do
    {:ok, color} = ColorManager.assign_color(user_id)

    socket =
      socket
      |> assign(:user_id, user_id)
      |> assign(:username, username)
      |> assign(user_color: %{user_id => color})

    send(self(), :after_join)
    {:ok, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    {:ok, _} =
      Presence.track(socket, socket.assigns.user_id, %{
        online_at: inspect(System.system_time(:second)),
        metadata: %{"source" => "customerOS"},
        username: socket.assigns.username,
        color: Map.get(socket.assigns.user_color, socket.assigns.user_id)
      })

    push(socket, "presence_state", Presence.list(socket))
    {:noreply, socket}
  end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  @impl true
  def handle_in("ping", payload, socket) do
    {:reply, {:ok, payload}, socket}
  end

  @impl true
  def terminate(_, socket) do
    Logger.info("User #{socket.assigns.user_id} left the channel")
    {:ok, socket}
  end

  # Add authorization logic here as required.
  # defp authorized?(_payload) do
  #   true
  # end
end
