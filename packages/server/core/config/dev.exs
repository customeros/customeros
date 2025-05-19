import Config

# Web endpoint configuration for development
config :core, Web.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4000],
  check_origin: false,
  code_reloader: true,
  debug_errors: true,
  secret_key_base: "YhZZx2x6lsMGrDZRg8cpZiRf9cQf6UWo9hz9wyqAb/5Ym+sc0cIfW5PS8yHi8hB5",
  watchers: [],
  reloadable_compilers: [:gettext, :elixir],
  reloadable_apps: [:ui, :backend],
  live_reload: [
    patterns: [
      ~r"priv/static/(?!uploads/).*(js|css|png|jpeg|jpg|gif|svg)$",
      ~r"priv/gettext/.*(po)$",
      ~r"lib/realtime_web/(controllers|live|components)/.*(ex|heex)$"
    ]
  ]

# Development-specific configurations
config :core, :realtime,
  dev_routes: true,
  app_env: :dev

# Development logger configuration
config :logger, :console, format: "[$level] $message\n"

# Phoenix development configurations
config :phoenix, :stacktrace_depth, 20
config :phoenix, :plug_init_mode, :runtime
config :phoenix_live_view, :debug_heex_annotations, true
