defmodule CustomerOsRealtimeWeb.SettingsChannel do
  @moduledoc """
  This is the Channel that tracks Settings view.
  """
  require Logger
  use CustomerOsRealtimeWeb, :channel
  alias CustomerOsRealtimeWeb.Presence

  @impl true
  def join("settings:lobby", _payload, socket) do
    Logger.debug("Reached join handler in settings_view_channel.ex")
    send(self(), :after_join)
    # {:ok, assign(socket, :user_id)}
    {:ok, socket}

    # if authorized?(payload) do
    #   {:ok, socket}
    # else
    #   {:error, %{reason: "unauthorized"}}
    # end
  end

  @impl true
  def handle_info(:after_join, socket) do
    {:ok, _} =
      Presence.track(socket, socket.assigns.user_id, %{
        online_at: inspect(System.system_time(:second))
      })

    push(socket, "presence_state", Presence.list(socket))
    # push(socket, "feed", %{list: feed_items(socket)})
    {:noreply, socket}
  end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  @impl true
  def handle_in("ping", payload, socket) do
    Logger.info("Reached ping handler in settings_view_channel.ex")
    {:reply, {:ok, payload}, socket}
  end

  # Add authorization logic here as required.
  # defp authorized?(_payload) do
  #   true
  # end
end
