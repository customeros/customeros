defmodule CustomerOsRealtimeWeb.OrganizationChannelTest do
  use CustomerOsRealtimeWeb.ChannelCase
  require Logger
  # alias CustomerOsRealtimeWeb.UserSocket

  setup do
    token = Phoenix.Token.sign(@endpoint, "user", "user.id")

    {:ok, _, socket} =
      CustomerOsRealtimeWeb.UserSocket
      # |> connect(%{"user_token" => token})
      # TODO(@max-openline): replace this direct assign to use `payload` during join action
      |> socket("user_id", %{user_id: "USER.ID", username: "Max Mustermann", typing: true})
      |> subscribe_and_join(CustomerOsRealtimeWeb.OrganizationChannel, "organization:lobby", %{
        "user_token" => token
      })

    # token = Phoenix.Token.sign(@endpoint, "user", "user.id")
    # {:ok, socket} = connect(UserSocket, %{"user_token" => token})
    # {:ok, _, ^socket} = subscribe_and_join(socket, "organization:lobby")

    # room: "organization:lobby"
    %{socket: socket}
  end

  test "ping replies with status ok", %{socket: socket} do
    ref = push(socket, "ping", %{"hello" => "there"})
    Logger.debug("Ref: #{inspect(ref)}")
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

  test "broadcasting presence", %{socket: socket} do
    {:ok, _, _} =
      subscribe_and_join(socket, "organization:lobby", %{})

    user_data = %{
      "USER.ID" => %{
        metas: [
          %{
            metadata: %{
              "source" => "customerOS"
            },
            online_at: inspect(System.system_time(:second)),
            username: "Max Mustermann",
            typing: true,
            phx_ref: "F7I72QkXPlaBcQKF"
          }
        ]
      }
    }

    diff_payload = %{
      leaves: %{},
      joins: %{
        "USER.ID" => %{
          metas: [
            %{
              metadata: %{
                "source" => "customerOS"
              },
              phx_ref: "F7I0nelpA_5MCgBk",
              online_at: inspect(System.system_time(:second)),
              username: "Max Mustermann",
              typing: true
            }
          ]
        }
      }
    }

    assert_push "presence_state", user_data
    assert_broadcast "presence_diff", diff_payload

    on_exit(fn ->
      for pid <- CustomerOsRealtimeWeb.Presence.fetchers_pids() do
        ref = Process.monitor(pid)
        assert_receive {:DOWN, ^ref, _, _, _}, 1000
      end
    end)
  end
end
