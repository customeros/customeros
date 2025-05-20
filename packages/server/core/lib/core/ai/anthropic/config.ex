defmodule Core.Ai.Anthropic.Config do
  @type t :: %__MODULE__{
          api_path: String.t(),
          api_key: String.t(),
          timeout: integer()
        }

  defstruct [
    :api_path,
    :api_key,
    ## 45 seconds
    timeout: 45_000
  ]

  def from_application_env do
    api_path =
      Application.get_env(:ai, :anthropic_api_path, "https://api.anthropic.com/v1/messages")

    api_key = Application.get_env(:ai, :anthropic_api_key)
    timeout = Application.get_env(:ai, :default_llm_timeout, 45_000)

    %__MODULE__{
      api_path: api_path,
      api_key: api_key,
      timeout: timeout
    }
  end

  def validate(%__MODULE__{} = config) do
    cond do
      is_nil(config.api_key) or config.api_key == "" ->
        {:error, "API key is required"}

      is_nil(config.api_path) or config.api_path == "" ->
        {:error, "API path is required"}

      true ->
        :ok
    end
  end
end
