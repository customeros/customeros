defmodule Web.Plugs.CORSWebSocket do
  import Plug.Conn

  def init(opts), do: opts

  def call(conn, opts) do
    if websocket_request?(conn) do
      # Skip CORS for WebSockets
      conn
    else
      CORSPlug.call(conn, CORSPlug.init(opts))
    end
  end

  defp websocket_request?(conn) do
    conn |> get_req_header("upgrade") |> List.first() |> String.downcase() == "websocket"
  end
end
