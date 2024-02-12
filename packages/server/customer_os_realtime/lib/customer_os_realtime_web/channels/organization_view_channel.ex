defmodule CustomerOsRealtimeWeb.OrganizationChannel do
  require Logger
  use CustomerOsRealtimeWeb, :channel
  alias CustomerOsRealtimeWeb.Presence

  @impl true
  def join("organization:lobby", _payload, socket) do
    Logger.debug("Reached join handler in organization_view_channel.ex")
    # %{"user_token" => user_token} = payload
    # case Phoenix.Token.verify(socket, "user", user_token) do
    #   {:ok, user_id} ->
    #     {:ok, assign(socket, :user_id, user_id)}
    #   {:error, _} ->
    #     :error
    # end
    # TODO assign user_id, username, typing from the payload to the socket.assigns
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
        online_at: inspect(System.system_time(:second)),
        metadata: %{"source" => "customerOS"},
        username: socket.assigns.username,
        typing: socket.assigns.typing
      })

    push(socket, "presence_state", Presence.list(socket))
    # push(socket, "feed", %{list: feed_items(socket)}) # TODO: push CRDT state to the client
    {:noreply, socket}
  end

  # # Helper functions
  # defp merge_states(crdt_state, client_state) do
  #   Map.merge(client_state, crdt_state)
  # end
  # defp feed_items(socket) do
  #   # Use the socket to fetch the feed items
  #   # if client has data already in the socket, use that
  #   # client_state = if Map.has_key?(payload, "body") do
  #   #   Map.fetch!(payload, "body")
  #   # else
  #   #   %{"body" => ""}
  #   # end
  #   client_state = %{"body" => ""}
  #   # Fetch current CRDT state
  #   crdt_state = CustomerOsRealtimeWeb.GlobalConfig.get(:payload)
  #   Logger.debug("CRDT state: #{inspect(crdt_state)}")

  #   # Merge CRDT state with client state
  #   merged_state = merge_states(crdt_state, client_state)
  #   Logger.debug("Merged state: #{inspect(merged_state)}")

  #   # Push merged state back to the client
  #   push(socket, "shout", merged_state)
  # end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  @impl true
  def handle_in("ping", payload, socket) do
    Logger.info("Reached ping handler in organization_view_channel.ex")
    {:reply, {:ok, payload}, socket}
  end

  # It is also common to receive messages from the client and
  # broadcast to everyone in the current topic (chat:lobby).
  @impl true
  def handle_in("shout", payload, socket) do
    Logger.info("Reached shout handler in organization_view_channel.ex")
    broadcast!(socket, "shout", payload)
    {:noreply, socket}
  end

  # Add authorization logic here as required.
  # defp authorized?(_payload) do
  #   true
  # end
end
