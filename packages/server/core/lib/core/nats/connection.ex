defmodule Core.Nats.Connection do
  use Supervisor

  @nats_connection :gnat

  def start_link(_args) do
    Supervisor.start_link(__MODULE__, nil, name: __MODULE__)
  end

  @impl true
  def init(_) do
    config = Core.Nats.Config.from_env()

    case Core.Nats.Config.validate(config) do
      :ok ->
        children = [
          {
            Gnat.ConnectionSupervisor,
            %{
              name: @nats_connection,
              connection_settings: build_connection_settings(config)
            }
          }
        ]

        Supervisor.init(children, strategy: :one_for_one)

      {:error, reason} ->
        {:error, reason}
    end
  end

  def conn, do: @nats_connection

  defp build_connection_settings(config) do
    initial_settings = [
      %{
        host: config.nats_node_1,
        port: config.nats_port
      }
    ]

    with_node_2 =
      if config.environment == "production" && config.nats_node_2 &&
           config.nats_node_2 != config.nats_node_1 do
        initial_settings ++
          [
            %{
              host: config.nats_node_2,
              port: config.nats_port
            }
          ]
      else
        initial_settings
      end

    settings =
      if config.environment == "production" && config.nats_node_3 &&
           config.nats_node_3 != config.nats_node_1 && config.nats_node_3 != config.nats_node_2 do
        with_node_2 ++
          [
            %{
              host: config.nats_node_3,
              port: config.nats_port
            }
          ]
      else
        with_node_2
      end

    settings
  end
end
