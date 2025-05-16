defmodule Realtime.MixProject do
  use Mix.Project

  def project do
    [
      app: :realtime,
      version: "0.1.0",
      elixir: "~> 1.14",
      elixirc_paths: elixirc_paths(Mix.env()),
      start_permanent: Mix.env() == :prod,
      aliases: aliases(),
      deps: deps(),
      releases: [
        realtime: [
          applications: [opentelemetry: :temporary]
        ]
      ],
      dialyzer: [plt_add_deps: :apps_direct]
    ]
  end

  def application do
    [
      mod: {Realtime.Application, []},
      extra_applications: [:logger, :runtime_tools]
    ]
  end

  defp elixirc_paths(:test), do: ["lib", "test/support"]
  defp elixirc_paths(_), do: ["lib"]

  defp deps do
    [
      {:telemetry_metrics, "~> 0.6"},
      {:telemetry_poller, "~> 1.0"},
      {:jason, "~> 1.2"},
      {:plug_cowboy, "~> 2.5"},
      {:phoenix, "~> 1.7.20"},
      {:phoenix_live_reload, "~> 1.2", only: :dev},
      {:phoenix_live_view, "~> 1.0.5"},
      {:phoenix_live_dashboard, "~> 0.8.6"},
      {:dns_cluster, "~> 0.1.1"},
      {:bandit, "~> 1.2"},
      {:dialyxir, "~> 1.4", only: [:dev, :test], runtime: false},
      {:credo, "~> 1.6", only: [:dev, :test], runtime: false},
      {:delta, "~> 0.2.0"},
      {:castore, "~> 1.0"},
      {:amqp, "~> 4.0"},
      {:opentelemetry, "~> 1.5"},
      {:opentelemetry_api, "~> 1.4"},
      {:opentelemetry_exporter, "~> 1.8"},
      {:opentelemetry_phoenix, "~> 2.0"},
      {:opentelemetry_bandit, "~> 0.2"},
      {:y_ex, "~> 0.7"},
      {:phoenix_ecto, "~> 4.6"},
      {:ecto_sql, "~> 3.12"},
      {:postgrex, "~> 0.20"},
      {:ecto_psql_extras, "~> 0.8"},
      {:temp, "~> 0.4"},
      {:absinthe, "~> 1.7"},
      {:absinthe_plug, "~> 1.5"},
      {:cors_plug, "~> 3.0"},
      {:grpc, "~> 0.10.1"},
      {:protobuf_generate, "~> 0.1.3", runtime: false},
      {:cowlib, "~> 2.15", override: true},
      {:gnat, "~> 1.10"},
      {:jetstream, "~> 0.0.9"}
    ]
  end

  defp aliases do
    [
      setup: ["deps.get", "ecto.setup"],
      dev: ["compile.script", "phx.server"],
      clean: ["deps.clean --unused --unlock"],
      "ecto.setup": ["ecto.create", "ecto.migrate"],
      proto: ["proto.fetch", "proto.gen"]
    ]
  end
end
