defmodule CustomerOsRealtimeWeb.FinderChannel do
  @moduledoc """
  This is the Channel that tracks Finder view.
  """
  require Logger
  use CustomerOsRealtimeWeb, :channel
  alias CustomerOsRealtimeWeb.Presence

  @impl true
  def join("finder:lobby", _payload, socket) do
    Logger.debug("Reached join handler in finder_view_channel.ex")
    send(self(), :after_join)
    {:ok, socket}
  end

  @impl true
  def handle_info(:after_join, socket) do
    {:ok, _} =
      Presence.track(socket, socket.assigns.user_id, %{
        online_at: inspect(System.system_time(:second))
      })

    push(socket, "presence_state", Presence.list(socket))
    {:noreply, socket}
  end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  @impl true
  def handle_in("ping", payload, socket) do
    Logger.info("Reached ping handler in finder_view_channel.ex")
    {:reply, {:ok, payload}, socket}
  end

  # Add authorization logic here as required.
  # defp authorized?(_payload) do
  #   true
  # end
end
