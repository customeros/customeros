defmodule Nats.Config do
  @type t :: %__MODULE__{
          environment: String.t(),
          nats_node_1: String.t(),
          nats_node_2: String.t(),
          nats_node_3: String.t(),
          nats_port: String.t()
        }

  defstruct [
    :environment,
    :nats_node_1,
    :nats_node_2,
    :nats_node_3,
    :nats_port
  ]

  def from_env do
    environment = Application.get_env(:ai, :environment)
    nats_node_1 = Application.get_env(:ai, :nats_node_1)
    nats_node_2 = Application.get_env(:ai, :nats_node_2)
    nats_node_3 = Application.get_env(:ai, :nats_node_3)
    nats_port = Application.get_env(:ai, :nats_port)

    %__MODULE__{
      environment: environment,
      nats_node_1: nats_node_1,
      nats_node_2: nats_node_2,
      nats_node_3: nats_node_3,
      nats_port: nats_port
    }
  end

  def validate(%__MODULE__{} = config) do
    cond do
      is_nil(config.nats_node_1) or config.nats_node_1 == "" ->
        {:error, "Nats config required"}

      is_nil(config.nats_port) or config.nats_port == "" ->
        {:error, "Nats port not set"}

      true ->
        :ok
    end
  end
end
