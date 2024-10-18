defmodule Realtime.OrganizationEventSubscriber do
  @moduledoc false
  use GenServer
  require Logger

  alias Realtime.EventStoreClient
  alias RealtimeWeb.Endpoint

  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def init(state) do
    stream_name = "notifyRealtime-v2"

    Logger.info("Subscribed to EventStore")
    Logger.info("Listening on persistent subscription: 'notifyRealtime-v2'")

    {:ok, subscription} =
      Spear.connect_to_persistent_subscription(EventStoreClient, self(), :all, stream_name)

    {:ok, Map.put(state, :subscription, subscription)}
  end

  def handle_info(
        %Spear.Event{
          body: %{
            "entity" => entity,
            "entityId" => entity_id,
            "tenant" => tenant,
            "create" => create,
            "update" => update,
            "delete" => delete
          }
        } =
          event,
        %{subscription: subscription} = state
      ) do
    channel_topic =
      case entity do
        "ORGANIZATION" ->
          "Organizations:#{tenant}"

        "CONTACT" ->
          "Contacts:#{tenant}"

        "CONTRACT" ->
          "Contracts:#{tenant}"

        "OPPORTUNITY" ->
          "Opportunities:#{tenant}"

        "SERVICE_LINE_ITEM" ->
          "ContractLineItems:#{tenant}"

        _ ->
          IO.puts("Unknown entity: #{entity}")
          nil
      end

    action_type =
      cond do
        create -> "APPEND"
        update -> "INVALIDATE"
        delete -> "DELETE"
      end

    if channel_topic != nil do
      Endpoint.broadcast!(channel_topic, "sync_group_packet", %{
        action: action_type,
        ids: [entity_id]
      })

      IO.puts("Broadcasted V1_EVENT_COMPLETED:#{action_type} to #{channel_topic}")
    end

    Spear.ack(EventStoreClient, subscription, event)

    {:noreply, state}
  end

  def handle_info({:dropped, reason}, state) do
    IO.puts("Subscription dropped: #{inspect(reason)}")
    {:stop, :normal, state}
  end
end
