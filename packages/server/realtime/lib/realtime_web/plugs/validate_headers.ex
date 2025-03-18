defmodule RealtimeWeb.Plugs.ValidateHeaders do
  import Plug.Conn

  @required_headers ["x-tenant"]

  def init(opts), do: opts

  def call(conn, _opts) do
    missing_headers =
      @required_headers
      |> Enum.filter(fn header -> get_req_header(conn, header) == [] end)

    if missing_headers == [] do
      context = %{tenant: get_req_header(conn, "x-tenant") |> List.first()}

      Absinthe.Plug.put_options(conn, context: context)
    else
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(
        400,
        Jason.encode!(%{error: "Missing required headers", missing: missing_headers})
      )
      |> halt()
    end
  end
end
