import Config

config :realtime,
  ecto_repos: [Realtime.Repo],
  generators: [timestamp_type: :utc_datetime]

config :realtime, RealtimeWeb.Endpoint,
  url: [host: "localhost"],
  adapter: Bandit.PhoenixAdapter,
  render_errors: [
    formats: [html: RealtimeWeb.ErrorHTML, json: RealtimeWeb.ErrorJSON],
    layout: false
  ],
  pubsub_server: Realtime.PubSub,
  live_view: [signing_salt: "jVLoUB9r"]

config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

config :phoenix, :json_library, Jason

config :opentelemetry, :processors,
  otel_batch_processor: %{
    exporter: {:otel_exporter_stdout, []}
  }

import_config "#{config_env()}.exs"
