defmodule CustomerOsRealtimeWeb.OrganizationChannelTest do
  use CustomerOsRealtimeWeb.ChannelCase
  alias CustomerOsRealtimeWeb.Conversation
  alias CustomerOsRealtimeWeb.UserSocket

  setup do
    # {:ok, _, socket} =
    #   CustomerOsRealtimeWeb.UserSocket
    #   |> connect(%{"user_token" => token})
    #   |> socket("user_id", %{some: :assign})
    #   |> subscribe_and_join(CustomerOsRealtimeWeb.OrganizationChannel, "organization:lobby")
    #   # # on_exit(fn ->
    #   # #   for pid <- CustomerOsRealtimeWeb.Presence.fetchers_pids() do
    #   # #     ref = Process.monitor(pid)
    #   # #     assert_receive {:DOWN, ^ref, _, _, _}, 1000
    #   # #   end
    #   # # end)
    token = Phoenix.Token.sign(@endpoint, "user", "user.id")
    {:ok, socket} = connect(UserSocket, %{"token" => token})
    {:ok, _, ^socket} = subscribe_and_join(
      socket, "organization:lobby", %{}
    )

    %{socket: socket} # room: "organization:lobby"
  end

  test "ping replies with status ok", %{socket: socket} do
    ref = push(socket, "ping", %{"hello" => "there"})
    assert_reply ref, :ok, %{"hello" => "there"}
  end

  test "shout broadcasts to organization:lobby", %{socket: socket} do
    push(socket, "shout", %{"hello" => "all"})
    assert_broadcast "shout", %{"hello" => "all"}
  end

  test "broadcasts are pushed to the client", %{socket: socket} do
    broadcast_from!(socket, "broadcast", %{"some" => "data"})
    assert_push "broadcast", %{"some" => "data"}
  end

  # test "broadcasting presence", %{socket: _socket, user: _user} do
  #   # {:ok, _, socket} = subscribe_and_join(
  #   #   socket, "room:#{room.id}", %{}
  #   # )
  #   _user_data = %{
  #     typing: false, user_id: user.id, username: user.username
  #   }
  #   assert_push "presence_state", user_data
  #   # assert_broadcast "presence_diff", user_data
  # end
end
