defmodule Web.DocumentController do
  use Web, :controller
  require Logger
  alias Core.Realtime.Documents

  def create(
        conn,
        %{
          "name" => _name,
          "body" => _body,
          "userId" => _user_id,
          "tenant" => _tenant,
          "icon" => _icon,
          "color" => _color,
          "lexicalState" => _lexical_state,
          "organizationId" => _organization_id
        } = params
      ) do
    handle_create(conn, params)
  end

  def create(conn, _params) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(400, Jason.encode!(%{error: "Invalid request body"}))
  end

  def index(conn, %{"organization_id" => organization_id}) do
    tenant = get_req_header(conn, "x-tenant") |> List.first()
    documents = Documents.list_by_organization(organization_id, tenant)

    json_response =
      documents
      |> Enum.map(fn doc ->
        if is_struct(doc), do: Map.from_struct(doc), else: doc
      end)
      |> Enum.map(&Core.Realtime.Util.to_camel_case_map/1)
      |> Jason.encode!()

    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, json_response)
  end

  defp handle_create(conn, params) do
    case Documents.create_document(params, parseDto: true) do
      {:ok, %{document: document}} ->
        json_response =
          document
          |> Map.from_struct()
          |> Core.Realtime.Util.to_camel_case_map()
          |> Jason.encode!()

        conn
        |> put_resp_content_type("application/json")
        |> send_resp(201, json_response)

      {:error, changeset} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(
          422,
          Jason.encode!(%{
            error: "Unprocessable entity",
            details: errors_from_changeset(changeset)
          })
        )
    end
  end

  defp errors_from_changeset(changeset) do
    Ecto.Changeset.traverse_errors(changeset, fn {msg, opts} ->
      Enum.reduce(opts, msg, fn {key, value}, acc ->
        String.replace(acc, "%{#{key}}", to_string(value))
      end)
    end)
  end
end
