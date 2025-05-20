defmodule Core.MixProject do
  use Mix.Project

  def project do
    [
      app: :core,
      version: "0.1.0",
      elixir: "~> 1.18",
      elixirc_paths: elixirc_paths(Mix.env()),
      start_permanent: Mix.env() == :prod,
      aliases: aliases(),
      deps: deps(),
      releases: [
        core: [
          applications: [opentelemetry: :temporary],
          include_executables_for: [:unix],
          strip_beams: Mix.env() == :prod,
        ]
      ],
      dialyzer: [plt_add_deps: :apps_direct]
    ]
  end

  def application do
    [
      mod: {Core.Application, []},
      extra_applications: [:logger, :runtime_tools]
    ]
  end

  defp elixirc_paths(:test), do: ["lib", "test/support"]
  defp elixirc_paths(_), do: ["lib"]

  defp deps do
    [
      {:absinthe, "~> 1.7"},
      {:absinthe_plug, "~> 1.5"},
      {:amqp, "~> 4.0"},
      {:bandit, "~> 1.2"},
      {:castore, "~> 1.0"},
      {:cors_plug, "~> 3.0"},
      {:cowlib, "~> 2.15", override: true},
      {:credo, "~> 1.7", only: [:dev, :test], runtime: false},
      {:delta, "~> 0.2.0"},
      {:dialyxir, "~> 1.4", only: [:dev, :test], runtime: false},
      {:dns_cluster, "~> 0.1.1"},
      {:ecto_psql_extras, "~> 0.8"},
      {:ecto_sql, "~> 3.12"},
      {:gnat, "~> 1.10"},
      {:grpc, "~> 0.10.1"},
      {:hackney, "~> 1.18"},
      {:jason, "~> 1.4"},
      {:jetstream, "~> 0.0.9"},
      {:opentelemetry, "~> 1.5"},
      {:opentelemetry_api, "~> 1.4"},
      {:opentelemetry_bandit, "~> 0.2"},
      {:opentelemetry_exporter, "~> 1.8"},
      {:opentelemetry_phoenix, "~> 2.0"},
      {:phoenix, "~> 1.7.20"},
      {:phoenix_ecto, "~> 4.6"},
      {:phoenix_live_dashboard, "~> 0.8.6"},
      {:phoenix_live_reload, "~> 1.2", only: :dev},
      {:phoenix_live_view, "~> 1.0.5"},
      {:plug_cowboy, "~> 2.5"},
      {:postgrex, "~> 0.20"},
      {:protobuf_generate, "~> 0.1.3", runtime: false},
      {:telemetry_metrics, "~> 0.6"},
      {:telemetry_poller, "~> 1.0"},
      {:temp, "~> 0.4"},
      {:y_ex, "~> 0.7"}
    ]
  end

  defp aliases do
    [
      clean: ["deps.clean --unused --unlock"],
      dev: ["compile.script", "phx.server"],
      "ecto.setup": ["ecto.create", "ecto.migrate"],
      proto: ["proto.fetch", "proto.gen"],
      setup: ["deps.get", "ecto.setup"],
      tidy: ["deps.get"]
    ]
  end
end
