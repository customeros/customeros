import Config

config :realtime, RealtimeWeb.Endpoint,
  cache_static_manifest: "priv/static/cache_manifest.json",
  check_origin: [
    "https://app.customeros.dev",
    "https://app.customeros.ai",
    "https://app.customeros.local",
    "https://frontera.customeros.ai",
    "https://frontera.openline.dev",
    "//*.localcan.dev"
  ]

config :logger, level: :info

config :realtime, :app_env, :prod
