defmodule Realtime.Nats.Consumer do
  use Jetstream.PullConsumer
  alias Gnat.Jetstream.API.{Consumer}

  def start_link(opts) do
    Jetstream.PullConsumer.start_link(__MODULE__, opts)
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
end
