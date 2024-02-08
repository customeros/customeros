defmodule CustomerOsRealtimeWeb.OrganizationChannelTest do
  use CustomerOsRealtimeWeb.ChannelCase

  setup do
    {:ok, _, socket} =
      RealtimeWeb.UserSocket
      |> socket("user_id", %{some: :assign})
      |> subscribe_and_join(CustomerOsRealtimeWeb.OrganizationChannel, "organization:lobby")
      # on_exit(fn ->
      #   for pid <- CustomerOsRealtimeWeb.Presence.fetchers_pids() do
      #     ref = Process.monitor(pid)
      #     assert_receive {:DOWN, ^ref, _, _, _}, 1000
      #   end
      # end)

    %{socket: socket}
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
end
