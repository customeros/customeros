defmodule Core.Nats.StreamSetup do
  use GenServer
  require Logger

  def start_link(_) do
    GenServer.start_link(__MODULE__, nil, name: __MODULE__)
  end

  @impl true
  def init(_) do
    # Schedule stream setup after initialization
    Process.send_after(self(), :setup_streams, 1000)
    {:ok, %{ready: false}}
  end

  @impl true
  def handle_info(:setup_streams, state) do
    if state.ready do
      {:noreply, state}
    else
      # Define streams to create
      streams = [
        {"ai", ["ai.>"]}
      ]

      # Try to set up all streams
      results =
        Enum.map(streams, fn {stream, subjects} ->
          {stream, Core.Nats.Streams.ensure_stream(Core.Nats.Connection.conn(), stream, subjects)}
        end)

      # Check if any stream setup failed
      failures = Enum.filter(results, fn {_, result} -> result != :ok end)

      if Enum.empty?(failures) do
        Logger.info("All NATS streams set up successfully")
        {:noreply, %{state | ready: true}}
      else
        # Log failures and retry
        Enum.each(failures, fn {stream, error} ->
          Logger.error("Failed to set up stream '#{stream}': #{inspect(error)}")
        end)

        # Retry after delay
        Process.send_after(self(), :setup_streams, 5000)
        {:noreply, state}
      end
    end
  end

  # Public API to check if streams are ready
  def ready? do
    GenServer.call(__MODULE__, :ready?)
  end

  @impl true
  def handle_call(:ready?, _from, state) do
    {:reply, state.ready, state}
  end
end
