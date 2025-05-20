defmodule Core.Ai.Anthropic.Service do
  alias Core.Ai.Anthropic.Models.Message
  alias Core.Ai.Anthropic.Models.AnthropicApiRequest
  alias Core.Ai.Anthropic.Config
  alias Core.Ai.Anthropic.Models.AskAIRequest

  require Logger

  @api_key_header "x-api-key"
  @anthropic_api_header "anthropic-version"
  @content_type_header "content-type"
  @default_api_version "2023-06-01"

  @spec ask(AskAIRequest.t(), Config.t()) :: {:ok, String.t()} | {:error, any()}
  def ask(%AskAIRequest{} = message, %Config{} = config) do
    with :ok <- AskAIRequest.validate(message),
         :ok <- Config.validate(config),
         {:ok, req_body} <- build_request(message),
         {:ok, response} <- execute(req_body, config) do
      {:ok, response}
    else
      {:error, reason} -> {:error, reason}
    end
  end

  defp build_request(%AskAIRequest{} = message) do
    model_name =
      case message.model do
        :claude_haiku -> "claude-3-5-haiku-20241022"
        :claude_sonnet -> "claude-3-5-sonnet-20241022"
      end

    request = %AnthropicApiRequest{
      model: model_name,
      messages: [%Message{role: "user", content: message.prompt}],
      max_tokens: message.max_output_tokens,
      temperature: message.model_temperature
    }

    request =
      if message.system_prompt && String.trim(message.system_prompt) != "" do
        %{request | system: message.system_prompt}
      else
        request
      end

    {:ok, request}
  end

  defp execute(req_body, config) do
    headers = [
      {@content_type_header, "application/json"},
      {@api_key_header, config.api_key},
      {@anthropic_api_header, @default_api_version}
    ]

    options = [
      timeout: config.timeout,
      recv_timeout: config.timeout
    ]

    case Jason.encode(req_body) do
      {:ok, json_body} ->
        case :hackney.request(:post, config.api_path, headers, json_body, options) do
          {:ok, status, _headers, client_ref} ->
            case :hackney.body(client_ref) do
              {:ok, body} -> process_response(status, body)
              {:error, reason} -> {:error, {:http_error, reason}}
            end

          {:error, reason} ->
            {:error, {:http_error, reason}}
        end

      {:error, reason} ->
        {:error, {:json_encode_error, reason}}
    end
  end

  defp process_response(status, body) when status in 200..299 do
    with {:ok, decoded} <- Jason.decode(body),
         content when is_list(content) <- Map.get(decoded, "content"),
         %{"type" => "text", "text" => text} <-
           Enum.find(content, &(Map.get(&1, "type") == "text")) do
      {:ok, String.trim(text)}
    else
      nil ->
        {:error, {:invalid_response, "No text content found in response"}}

      _ ->
        {:error, {:invalid_response, "Invalid response format"}}
    end
  end

  defp process_response(status, body) do
    case Jason.decode(body, keys: :atoms) do
      {:ok, %{error: %{type: type, message: message}}} ->
        {:error, {:api_error, "#{type}: #{message}"}}

      _ ->
        {:error, {:http_error, "API request failed with status #{status}: #{body}"}}
    end
  end
end
