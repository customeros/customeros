defmodule CustomerOsRealtime.ChannelClient do
  use GenServer
  require Logger

  # alias Phoenix.Socket.Broadcast
  alias CustomerOsRealtimeWeb.Endpoint

  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def sync_group_packet(channel, message) do
    GenServer.cast(__MODULE__, {:sync_group_packet, channel, message})
  end

  def init(state) do
    {:ok, state}
  end

  def handle_cast({:sync_group_packet, channel, message}, state) do
    sub = Endpoint.subscribe(channel)

    case sub do
      :ok ->
        IO.puts("Subscribed to channel #{channel}")

      {:error, _} ->
        IO.puts("Error subscribing to channel #{channel}")
    end

    result =
      Endpoint.broadcast_from(self(), channel, "sync_group_packet", %{
        payload: message
      })

    case result do
      :ok ->
        IO.puts("Broadcasting message on channel #{channel}")
        {:noreply, state}

      {:error, _} ->
        IO.puts("Broadcast error on channel #{channel}")
        {:noreply, state}
    end
  end
end
