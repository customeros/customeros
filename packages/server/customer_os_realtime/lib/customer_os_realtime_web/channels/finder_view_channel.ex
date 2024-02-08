defmodule CustomerOsRealtimeWeb.FinderChannel do
  require Logger
  use CustomerOsRealtimeWeb, :channel

  @impl true
  def join("finder:lobby", _payload, socket) do
    Logger.debug("Reached join handler in finder_view_channel.ex")
    {:ok, socket}

    # if authorized?(payload) do
    #   {:ok, socket}
    # else
    #   {:error, %{reason: "unauthorized"}}
    # end
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
