import Config

# Helper functions for environment variables
get_env = fn key, default -> System.get_env(key, default) end
get_env_integer = fn key, default -> String.to_integer(get_env.(key, default)) end

get_env_boolean = fn key, default ->
  case get_env.(key, default) do
    val when val in ["true", "1"] -> true
    _ -> false
  end
end

# Load environment variables from .env in development
if config_env() == :dev do
  case File.read(".env") do
    {:ok, content} ->
      content
      |> String.split("\n")
      |> Enum.map(&String.trim/1)
      |> Enum.filter(&(&1 != "" && !String.starts_with?(&1, "#")))
      |> Enum.each(fn line ->
        [key, value] = String.split(line, "=", parts: 2)
        System.put_env(key, value)
      end)

    _ ->
      :ok
  end
end

# Server configuration
if get_env.("PHX_SERVER", nil), do: config(:realtime, Web.Endpoint, server: true)

# Database configuration
config :core, Core.Repo,
  username: get_env.("POSTGRES_USER", "postgres"),
  password: get_env.("POSTGRES_PASSWORD", "password"),
  hostname: get_env.("POSTGRES_HOST", "localhost"),
  database: get_env.("POSTGRES_DB", "openline"),
  port: get_env_integer.("POSTGRES_PORT", "5555"),
  stacktrace: true,
  show_sensitive_data_on_connection_error: true,
  pool_size: 10

# OpenTelemetry configuration
config :opentelemetry, :resource, service: %{name: "Realtime"}

config :opentelemetry, :processors,
  otel_batch_processor: %{
    exporter: {
      :opentelemetry_exporter,
      %{
        endpoints: [
          {:http, get_env.("JAEGER_AGENT_HOST", "localhost"),
           get_env_integer.("JAEGER_AGENT_PORT", "4318"), []}
        ]
      }
    }
  }

# Nats configuration
if config_env() == :prod do
  # Set NATS environment to production in production mode
  config :core, :nats, environment: "production"
end

# Configure NATS hosts if environment variables are set
Enum.each(1..3, fn node ->
  host_var = "NATS_HOST_NODE_#{node}"

  if host = get_env.(host_var, "localhost") do
    if host != "", do: config(:core, :nats, String.to_atom("nats_node_#{node}"))
  end
end)

# Configure NATS port if environment variable is set
if port_str = get_env.("NATS_PORT", "4222") do
  if port_str != "" do
    config :core, :nats, nats_port: get_env_integer.("NATS_PORT", "4222")
  end
end

# AI configuration
config :core, :ai,
  anthropic_api_path: "https://api.anthropic.com/v1/messages",
  anthropic_api_key: get_env.("ANTHROPIC_API_KEY", nil)

# Production environment specific configuration
if config_env() == :prod do
  secret_key_base =
    get_env.("SECRET_KEY_BASE", nil) ||
      raise """
      environment variable SECRET_KEY_BASE is missing.
      You can generate one by calling: mix phx.gen.secret
      """

  host = get_env.("PHX_HOST", "example.com")
  port = get_env_integer.("PORT", "4000")

  config :core, :dns_cluster_query, get_env.("DNS_CLUSTER_QUERY", nil)

  config :core, Web.Endpoint,
    url: [host: host, port: 443, scheme: "https"],
    http: [
      ip: {0, 0, 0, 0, 0, 0, 0, 0},
      port: port
    ],
    secret_key_base: secret_key_base
end
