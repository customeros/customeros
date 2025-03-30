import Config

if System.get_env("PHX_SERVER") do
  config :realtime, RealtimeWeb.Endpoint, server: true
end

if config_env() == :dev do
  ".env"
  |> File.stream!()
  |> Stream.map(&String.trim/1)
  |> Stream.filter(&(&1 != "" && !String.starts_with?(&1, "#")))
  |> Enum.each(fn line ->
    [key, value] = String.split(line, "=", parts: 2)
    System.put_env(key, value)
  end)
end

jeager_host = System.get_env("JAEGER_AGENT_HOST", "localhost")
jeager_port = String.to_integer(System.get_env("JAEGER_AGENT_PORT", "4318"))
postgres_host = System.get_env("POSTGRES_HOST", "localhost")
postgres_port = String.to_integer(System.get_env("POSTGRES_PORT", "5432"))
postgres_username = System.get_env("POSTGRES_USER", "postgres")
postgres_password = System.get_env("POSTGRES_PASSWORD", "password")
postgres_database = System.get_env("POSTGRES_DB", "realtime")

config :realtime, Realtime.Repo,
  username: postgres_username,
  password: postgres_password,
  hostname: postgres_host,
  database: postgres_database,
  port: postgres_port,
  stacktrace: true,
  show_sensitive_data_on_connection_error: true,
  pool_size: 10

config :opentelemetry,
       :resource,
       service: %{name: "Realtime"}

config :opentelemetry,
       :processors,
       otel_batch_processor: %{
         exporter:
           {:opentelemetry_exporter, %{endpoints: [{:http, jeager_host, jeager_port, []}]}}
       }

if config_env() == :prod do
  secret_key_base =
    System.get_env("SECRET_KEY_BASE") ||
      raise """
      environment variable SECRET_KEY_BASE is missing.
      You can generate one by calling: mix phx.gen.secret
      """

  host = System.get_env("PHX_HOST") || "example.com"
  port = String.to_integer(System.get_env("PORT") || "4000")

  config :realtime, :dns_cluster_query, System.get_env("DNS_CLUSTER_QUERY")

  config :realtime, RealtimeWeb.Endpoint,
    url: [host: host, port: 443, scheme: "https"],
    http: [
      ip: {0, 0, 0, 0, 0, 0, 0, 0},
      port: port
    ],
    secret_key_base: secret_key_base
end
