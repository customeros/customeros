import Config

# Core application configuration
config :core,
  ecto_repos: [Core.Realtime.Repo],
  generators: [timestamp_type: :utc_datetime]

# Web endpoint configuration
config :core, Web.Endpoint,
  url: [host: "localhost"],
  adapter: Bandit.PhoenixAdapter,
  render_errors: [
    formats: [html: Web.ErrorHTML, json: Web.ErrorJSON],
    layout: false
  ],
  pubsub_server: Core.Realtime.PubSub,
  live_view: [signing_salt: "jVLoUB9r"]

# Logger configuration
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Phoenix configuration
config :phoenix, :json_library, Jason

# Default OpenTelemetry configuration (overridden in runtime.exs)
config :opentelemetry, :processors,
  otel_batch_processor: %{
    exporter: {:otel_exporter_stdout, []}
  }

# Nats config
config :core, :nats,
  environment: "dev",
  nats_node_1: "localhost",
  nats_node_2: "localhost",
  nats_node_3: "localhost",
  nats_port: 4222

# AI configuration
config :core, :ai,
  anthropic_api_path: "https://api.anthropic.com/v1/messages",
  default_llm_timeout: 45_000

config :core, :rabbitmq,
  url: System.get_env("RABBITMQ_URL") || "amqp://guest:guest@localhost:5672"

# Import environment specific config
import_config "#{config_env()}.exs"
