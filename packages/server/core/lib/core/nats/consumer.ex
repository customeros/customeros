defmodule Core.Nats.Consumer do
  use Jetstream.PullConsumer
  require Logger
  alias Gnat.Jetstream.API.{Consumer}

  def start_link(opts) do
    case ensure_connection() do
      :ok ->
        Jetstream.PullConsumer.start_link(__MODULE__, opts)

      {:error, reason} ->
        Logger.error("Failed to start NATS consumer: #{inspect(reason)}")
        {:error, reason}
    end
  end

  @impl true
  def init(opts) do
    stream = Keyword.fetch!(opts, :stream)
    consumer = Keyword.fetch!(opts, :consumer)
    deliver_group = Keyword.fetch!(opts, :deliver_group)
    filter_subject = Keyword.fetch!(opts, :filter_subject)
    handler = Keyword.fetch!(opts, :handler)
    conn = :gnat

    case Consumer.create(conn, %Consumer{
           durable_name: consumer,
           stream_name: stream,
           deliver_group: deliver_group,
           filter_subject: filter_subject,
           deliver_policy: :all
         }) do
      {:ok, _resp} ->
        IO.puts("✅ Created or found consumer #{stream}.#{consumer}")

      {:error, %{error: %{"code" => 409}}} ->
        # Consumer already exists — that's okay
        IO.puts("ℹ️ Consumer #{stream}.#{consumer} already exists")

      {:error, reason} ->
        IO.warn("⚠️ Failed to create consumer #{stream}.#{consumer}: #{inspect(reason)}")
    end

    {:ok, %{stream: stream, consumer: consumer, handler: handler},
     connection_name: :gnat, stream_name: stream, consumer_name: consumer}
  end

  @impl true
  def handle_message(message, state) do
    state.handler.handle_message(message)
    {:ack, state}
  end

  defp ensure_connection(attempts \\ 30, delay \\ 1000) do
    case Process.whereis(:gnat) do
      nil ->
        if attempts > 0 do
          Logger.debug("Waiting for NATS connection to start (#{attempts} attempts left)", [])
          Process.sleep(delay)
          ensure_connection(attempts - 1, delay)
        else
          {:error, :connection_not_available}
        end

      pid ->
        # Simple check using basic NATS operations
        test_topic = "ai.connection.test.#{System.system_time()}"
        test_payload = "ping"

        try do
          # Basic pub operation - if this succeeds, connection is working
          case Gnat.pub(pid, test_topic, test_payload) do
            :ok ->
              Logger.debug("NATS connection verified with successful publish", [])
              :ok
          end
        rescue
          e ->
            Logger.warning("Error in NATS connection check: #{inspect(e)}", [])

            if attempts > 0 do
              Process.sleep(delay)
              ensure_connection(attempts - 1, delay)
            else
              {:error, e}
            end
        end
    end
  end
end
