defmodule Core.Nats.Streams do
  require Logger
  alias Gnat.Jetstream.API.Stream

  @doc """
  Ensures a stream exists, creating it if necessary.
  Returns :ok on success or {:error, reason} on failure.
  """
  def ensure_stream(conn_name, stream_name, subjects) do
    case stream_exists?(conn_name, stream_name) do
      true ->
        Logger.info("Stream '#{stream_name}' exists")
        :ok

      false ->
        Logger.info("Creating stream '#{stream_name}'")
        create_stream(conn_name, stream_name, subjects)
    end
  end

  defp stream_exists?(conn_name, stream_name) do
    case Stream.info(conn_name, stream_name) do
      {:ok, _stream_info} ->
        true

      {:error, %{status: 404}} ->
        false

      {:error, reason} ->
        Logger.error("Error checking if stream exists: #{inspect(reason)}")
        false
    end
  end

  defp create_stream(conn_name, stream_name, subjects) when is_list(subjects) do
    env = Application.get_env(:ai, :environment, "dev")
    replicas = get_replicas(env)

    stream_config = %Stream{
      name: stream_name,
      subjects: subjects,
      # retention can be :limits, :interest, or :workqueue
      retention: :workqueue,
      max_consumers: -1,
      max_msgs: -1,
      max_bytes: -1,
      max_age: 0,
      max_msgs_per_subject: -1,
      max_msg_size: -1,
      storage: :file,
      num_replicas: replicas
    }

    case Stream.create(conn_name, stream_config) do
      {:ok, _response} ->
        Logger.info("Successfully created stream '#{stream_name}'")
        :ok

      {:error, %{status: 400, error: %{"code" => 400, "description" => desc}}} ->
        if String.contains?(desc, "already exists") do
          Logger.error("Stream '#{stream_name}' already exists")
          :ok
        else
          Logger.error("Error creating stream: #{desc}")
          {:error, desc}
        end

      {:error, reason} ->
        Logger.error("Failed to create stream '#{stream_name}': #{inspect(reason)}")
        {:error, reason}
    end
  end

  defp get_replicas("production"), do: 3
  defp get_replicas(_), do: 1
end
