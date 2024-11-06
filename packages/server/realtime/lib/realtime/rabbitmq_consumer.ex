defmodule Realtime.RabbitMQConsumer do
  use GenServer
  require Logger
  require Jason
  alias AMQP.{Connection, Channel}
  alias RealtimeWeb.Endpoint

  @moduledoc false

  @queue_name "notifications"
  @entityToChannelMap %{
    "ORGANIZATION" => "Organizations",
    "CONTACT" => "Contacts",
    "CONTRACT" => "Contracts",
    "OPPORTUNITY" => "Opportunities",
    "SERVICE_LINE_ITEM" => "ContractLineItems",
    "FLOW" => "Flows",
    "FLOW_PARTICIPANT" => "FlowContacts",
    "SYSTEM" => "System"
  }

  def start_link(_) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def init(_) do
    rabbitmq_url = System.get_env("RABBITMQ_URL")
    {:ok, conn} = Connection.open(rabbitmq_url)
    {:ok, channel} = Channel.open(conn)

    case AMQP.Basic.consume(channel, "notifications") do
      {:ok, _} -> Logger.info("Consuming queue: #{@queue_name}")
      _ -> Logger.error("Cannot consume queue: #{@queue_name}")
    end

    {:ok, %{channel: channel}}
  end

  def handle_info({:basic_deliver, payload, meta}, state) do
    Logger.info("Received message on queue: #{@queue_name}")

    case Jason.decode(payload) do
      {:ok, parsed} ->
        %{
          "tenant" => tenant,
          "entityType" => entity_type,
          "entityIds" => entity_ids,
          "create" => create,
          "update" => update,
          "delete" => delete
        } = parsed

        message = Map.get(parsed, "message")
        channel_topic_prefix = Map.get(@entityToChannelMap, entity_type, :unknown)

        channel_topic =
          case channel_topic_prefix do
            :unknown ->
              Logger.warning("Unknown entity: #{entity_type}")
              nil

            value ->
              "#{value}:#{tenant}"
          end

        action_type =
          cond do
            create -> "APPEND"
            update -> "INVALIDATE"
            delete -> "DELETE"
            true -> message
          end

        case channel_topic do
          nil ->
            Logger.warning(
              "No channel_topic detected for entity:#{entity_type} - tenant:#{tenant}, will ack and do nothing."
            )

          _ ->
            Endpoint.broadcast!(channel_topic, "sync_group_packet", %{
              action: action_type,
              ids: entity_ids
            })

            Logger.info(
              "Broadcasted notification:#{action_type} to #{channel_topic} for #{tenant}"
            )
        end

      _ ->
        Logger.error("Failed decoding payload from JSON")
    end

    AMQP.Basic.ack(state.channel, meta.delivery_tag)

    {:noreply, state}
  end

  def handle_info({:basic_consume_ok, _meta}, state) do
    Logger.info("Successfully subscribed to the queue: #{@queue_name}")
    {:noreply, state}
  end

  def handle_info({:basic_cancel_ok, _meta}, state) do
    Logger.info("Subscription to the queue was cancelled.")
    {:noreply, state}
  end

  def handle_info({:basic_cancel, _meta}, state) do
    Logger.warning("Subscription to the queue was unexpectedly cancelled.")
    {:stop, :normal, state}
  end

  def handle_info(msg, state) do
    IO.puts("Unhandled message: #{inspect(msg)}")
    {:noreply, state}
  end
end
