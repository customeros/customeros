import Config

config :core, Web.Endpoint,
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

config :core, :app_env, :prod
