defmodule Web.Endpoint do
  use Phoenix.Endpoint, otp_app: :core

  # The session will be stored in the cookie and signed,
  # this means its contents can be read but not tampered with.
  # Set :encryption_salt if you would also like to encrypt it.
  @session_options [
    store: :cookie,
    key: "_customer_os_realtime_key",
    signing_salt: "m4WrTROE",
    same_site: "Lax"
  ]

  socket "/live", Phoenix.LiveView.Socket,
    websocket: [connect_info: [session: @session_options]],
    longpoll: [connect_info: [session: @session_options]]

  socket "/socket", Web.UserSocket, websocket: true, longpoll: false
  # Serve at "/" the static files from "priv/static" directory.
  #
  # You should set gzip to true if you are running phx.digest
  # when deploying your static files in production.
  plug Plug.Static,
    at: "/",
    from: :core,
    gzip: false,
    only: Web.static_paths()

  # Code reloading can be explicitly enabled under the
  # :code_reloader configuration of your endpoint.
  if code_reloading? do
    socket "/phoenix/live_reload/socket", Phoenix.LiveReloader.Socket
    plug Phoenix.LiveReloader
    plug Phoenix.CodeReloader
  end

  plug Phoenix.LiveDashboard.RequestLogger,
    param_key: "request_logger",
    cookie_key: "request_logger"

  # plug RealtimeWeb.Plugs.CORSWebSocket,
  #   origin: "*",
  #   methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
  #   headers: [
  #     "Authorization",
  #     "Content-Type",
  #     "Accept",
  #     "X-Tenant",
  #     "X-User-Id",
  #     "X-User-Email",
  #     "X-User-Roles",
  #     "X-User-Name",
  #     "X-Request-Id",
  #     "X-OPENLINE-API-KEY",
  #     "X-CUSTOMER-OS-API-KEY",
  #     "X-OPENLINE-USERNAME",
  #     "X-OPENLINE-TENANT"
  #   ],
  #   expose: [
  #     "Authorization"
  #   ],
  #   max_age: 86400,
  #   send_preflight_response?: true

  plug Plug.RequestId
  plug Plug.Telemetry, event_prefix: [:phoenix, :endpoint]

  plug Plug.Parsers,
    parsers: [:urlencoded, :multipart, :json],
    pass: ["*/*"],
    json_decoder: Phoenix.json_library()

  plug Plug.MethodOverride
  plug Plug.Head
  plug Plug.Session, @session_options
  plug Web.Router
end
