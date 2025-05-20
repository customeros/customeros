defmodule Core.Nats.Config do
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
    [{_, environment}, {_, nats_node_1}, {_, nats_node_2}, {_, nats_node_3}, {_, nats_port}] =
      Application.get_env(:core, :nats)

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
