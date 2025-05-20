defmodule Core.Application do
  @moduledoc false
  use Application

  @impl true
  def start(_type, _args) do
    OpentelemetryBandit.setup()
    OpentelemetryPhoenix.setup(adapter: :bandit)

    children = [
      Web.Telemetry,
      Core.Repo,
      {DNSCluster, query: Application.get_env(:core, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Realtime.PubSub},
      Web.Presence,
      Web.Endpoint,
      Core.Realtime.ColorManager,
      Core.Realtime.StoreManager
    ]

    env = Application.get_env(:core, :app_env, :prod)

    children =
      if env != :test do
        children ++
          [
            Core.Realtime.RabbitMQConsumer,
            Core.Nats.Supervisor
          ]
      else
        children
      end

    opts = [strategy: :one_for_one, name: Core.Realtime.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    Web.Endpoint.config_change(changed, removed)
    :ok
  end
end
