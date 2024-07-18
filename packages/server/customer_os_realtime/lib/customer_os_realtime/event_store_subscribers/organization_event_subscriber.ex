defmodule CustomerOsRealtime.OrganizationEventSubscriber do
  use GenServer

  alias CustomerOsRealtime.EventStoreClient
  alias CustomerOsRealtime.ChannelClient

  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def init(state) do
    stream_name = "organization-v1"

    {:ok, subscription} =
      Spear.connect_to_persistent_subscription(EventStoreClient, self(), :all, stream_name)

    {:ok, Map.put(state, :subscription, subscription)}
  end

  def handle_info(
        %Spear.Event{type: type, metadata: %{stream_name: stream_name} = _metadata} =
          event,
        %{subscription: subscription} = state
      ) do
    "organization-v1-" <> tenant_organization_id = stream_name
    [tenant, organization_id] = String.split(tenant_organization_id, "-")
    channel_topic = "Organizations:" <> tenant

    dbg(tenant)
    dbg(organization_id)

    case type do
      "V1_ORGANIZATION_CREATE" ->
        IO.puts(~c"Organization Created")

      "V1_ORGANIZATION_UPDATE" ->
        IO.puts(~c"Organization Updated")

        ChannelClient.sync_group_packet(channel_topic, %{
          action: "INVALIDATE",
          ids: [organization_id]
        })

      _ ->
        IO.puts("Unknown event type: #{type}")
    end

    Spear.ack(EventStoreClient, subscription, event)

    {:noreply, state}
  end

  def handle_info({:dropped, reason}, state) do
    IO.puts("Subscription dropped: #{inspect(reason)}")
    {:stop, :normal, state}
  end
end
