defmodule Core.Ai.Anthropic.Models do
  @moduledoc """
  Structs for representing Anthropic API requests and responses.
  """

  defmodule AskAIRequest do
    @derive Jason.Encoder
    @moduledoc """
    Represents a request to ask a question to Claude.
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
      model_temperature: 0.7
    ]

    @doc """
    Validates the request and returns errors if any.
    """
    def validate(%__MODULE__{} = request) do
      cond do
        request.model not in [:claude_haiku, :claude_sonnet] ->
          {:error, "model must be :claude_haiku or :claude_sonnet"}

        is_nil(request.prompt) or request.prompt == "" ->
          {:error, "prompt cannot be empty"}

        true ->
          :ok
      end
    end
  end

  defmodule AnthropicApiRequest do
    @derive Jason.Encoder
    @moduledoc """
    The JSON structure for a Claude API request.
    """

    @type t :: %__MODULE__{
            model: String.t(),
            system: String.t() | nil,
            messages: [Message.t()],
            max_tokens: integer() | nil,
            temperature: float() | nil
          }

    defstruct [:model, :system, :messages, :max_tokens, :temperature]
  end

  defmodule Message do
    @derive Jason.Encoder
    @moduledoc """
    A single message in a conversation.
    """

    @type t :: %__MODULE__{
            role: String.t(),
            content: String.t() | map()
          }

    defstruct [:role, :content]
  end

  defmodule Content do
    @derive Jason.Encoder
    @moduledoc """
    Content structure in the API response.
    """

    @type t :: %__MODULE__{
            type: String.t(),
            text: String.t()
          }

    defstruct [:type, :text]
  end

  defmodule Usage do
    @derive Jason.Encoder
    @moduledoc """
    Token usage information.
    """

    @type t :: %__MODULE__{
            input_tokens: integer(),
            output_tokens: integer()
          }

    defstruct [:input_tokens, :output_tokens]
  end

  defmodule AnthropicApiResponse do
    @derive Jason.Encoder
    @moduledoc """
    The JSON structure for a Claude API response.
    """

    @type t :: %__MODULE__{
            id: String.t(),
            type: String.t(),
            role: String.t(),
            content: [Content.t()],
            model: String.t(),
            stop_reason: String.t() | nil,
            usage: Usage.t()
          }

    defstruct [:id, :type, :role, :content, :model, :stop_reason, :usage]
  end

  defmodule ErrorResponse do
    @derive Jason.Encoder
    @moduledoc """
    Error response structure.
    """

    @type t :: %__MODULE__{
            error: %{
              type: String.t(),
              message: String.t()
            }
          }

    defstruct error: %{type: nil, message: nil}
  end
end
