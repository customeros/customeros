defmodule Core.Ai.AskAi.Service do
  defmodule AskAIRequest do
    @derive Jason.Encoder
    @moduledoc """
    Represents a request to ask a question to an AI Model.
    """

    @type model :: :claude_haiku | :claude_sonnet

    @type t :: %__MODULE__{
            model: model(),
            prompt: String.t(),
            system_prompt: String.t() | nil,
            max_output_tokens: integer() | nil,
            model_temperature: float() | nil
          }

    defstruct [
      :model,
      :prompt,
      :system_prompt,
      max_output_tokens: 1024,
      model_temperature: 0.2
    ]
  end

  @spec ask(AskAIRequest.t()) :: {:ok, String.t()} | {:error, any()}
  def ask(%AskAIRequest{} = message) do
    case message.model do
      model when model in [:claude_haiku, :claude_sonnet] ->
        anthropic_request = %Core.Ai.Anthropic.Models.AskAIRequest{
          model: message.model,
          prompt: message.prompt,
          system_prompt: message.system_prompt,
          max_output_tokens: message.max_output_tokens,
          model_temperature: message.model_temperature
        }

        Core.Ai.Anthropic.Service.ask(anthropic_request, Core.Ai.Anthropic.Config.from_application_env())

      unsupported_model ->
        {:error, {:unsupported_model, unsupported_model}}
    end
  end
end
