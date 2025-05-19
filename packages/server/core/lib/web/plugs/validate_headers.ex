defmodule Web.Plugs.ValidateHeaders do
  import Plug.Conn

  def init(opts), do: opts

  def call(conn, _opts) do
    tenant =
      case get_req_header(conn, "x-tenant") do
        [] -> get_req_header(conn, "x-openline-tenant") |> List.first()
        [value] -> value
      end

    username =
      case get_req_header(conn, "x-username") do
        [] -> get_req_header(conn, "x-openline-username") |> List.first()
        [value] -> value
      end

    # internal_api_key = get_req_header(conn, "x-openline-api-key") |> List.first()

    # api_token = System.get_env("API_TOKEN")

    # if internal_api_key && internal_api_key != api_token do
    #   conn
    #   |> put_resp_content_type("application/json")
    #   |> send_resp(
    #     401,
    #     Jason.encode!(%{error: "Unauthorized", message: "Invalid internal app token"})
    #   )
    #   |> halt()
    # else
    #   conn
    # end

    if !tenant do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(
        401,
        Jason.encode!(%{error: "Unauthorized", message: "missing tenant header"})
      )
      |> halt()
    end

    if !username do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(
        401,
        Jason.encode!(%{error: "Unauthorized", message: "missing username header"})
      )
      |> halt()
    end

    context = %{
      tenant: tenant,
      username: username
    }

    Absinthe.Plug.put_options(conn, context: context)
  end
end
