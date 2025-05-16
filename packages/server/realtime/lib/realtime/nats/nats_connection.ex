defmodule Realtime.Nats.Connection do
  use Supervisor

  @nats_connection :gnat

  def start_link(_args) do
    Supervisor.start_link(__MODULE__, nil, name: __MODULE__)
  end

  @impl true
  def init(_) do
    host = System.get_env("NATS_HOST")
    port = System.get_env("NATS_PORT") |> String.to_integer()

    children = [
      {
        Gnat.ConnectionSupervisor,
        %{
          name: @nats_connection,
          connection_settings: [
            %{
              host: host,
              port: port
            }
          ]
        }
      }
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end

  def conn, do: @nats_connection
end
