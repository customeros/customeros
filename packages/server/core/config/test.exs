import Config

config :core, Web.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4002],
  secret_key_base: "pcrqicRFh5xqnXMCq/W9kaZYhIJtXytsvf5L5Janxrk6VBFToY9Gr6Rjz+AaPhV+",
  server: false

config :core, Realtime.Mailer, adapter: Swoosh.Adapters.Test

config :swoosh, :api_client, false

config :logger, level: :warning

config :phoenix, :plug_init_mode, :runtime

config :core, :app_env, :test
