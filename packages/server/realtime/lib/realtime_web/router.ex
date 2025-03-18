defmodule RealtimeWeb.Router do
  use RealtimeWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :protect_from_forgery
    plug :put_secure_browser_headers
  end

  pipeline :api do
    plug :accepts, ["json"]
    plug RealtimeWeb.Plugs.ValidateHeaders
  end

  pipeline :graphql do
    plug RealtimeWeb.Plugs.ValidateHeaders
  end

  scope "/", RealtimeWeb do
    pipe_through :browser
  end

  scope "/graphql" do
    pipe_through :graphql
    forward "/", Absinthe.Plug, schema: RealtimeWeb.Graphql.Schema
  end

  forward "/graphiql",
          Absinthe.Plug.GraphiQL,
          schema: RealtimeWeb.Graphql.Schema,
          interface: :simple

  scope "/api", RealtimeWeb do
    pipe_through :api

    resources "/documents", DocumentController, only: [:create]
    post "/organizations/:organization_id/documents", DocumentController, :create
    get "/organizations/:organization_id/documents", DocumentController, :index
  end

  # Enable LiveDashboard and Swoosh mailbox preview in development
  if Application.compile_env(:realtime, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through [:fetch_session, :protect_from_forgery]

      live_dashboard "/dashboard", metrics: RealtimeWeb.Telemetry
      # forward "/mailbox", Plug.Swoosh.MailboxPreview
    end
  end
end
