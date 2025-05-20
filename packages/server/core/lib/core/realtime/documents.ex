defmodule Core.Realtime.Documents do
  @moduledoc """
  The Documents context.
  """
  import Ecto.Query, warn: false
  require Logger
  alias Core.Repo
  alias Core.Realtime.Documents.{Document, OrganizationDocument}

  def create_document(attrs \\ %{}, parseDto \\ false) do
    payload =
      case parseDto do
        true -> fromDto(attrs)
        _ -> attrs
      end

    organization_id = Map.get(payload, :organization_id)

    Ecto.Multi.new()
    |> Ecto.Multi.insert(:document, Document.changeset(%Document{}, payload))
    |> maybe_insert_organization_document(organization_id)
    |> maybe_initialize_ydoc(Map.get(payload, :lexical_state))
    |> Repo.transaction()
    |> case do
      {:ok, %{document: document} = result} ->
        decorated_doc = %Document{document | organization_id: organization_id}
        {:ok, Map.put(result, :document, decorated_doc)}

      error ->
        error
    end
  end

  def update_document(attrs \\ %{}) do
    with %Document{} = document <- Repo.get(Document, attrs.id) do
      document
      |> Document.update_changeset(attrs)
      |> Repo.update()
    else
      nil -> {:error, :not_found}
    end
  end

  def delete_document(id) do
    with %Document{} = document <- Repo.get(Document, id) do
      document
      |> Repo.delete()
    else
      nil -> {:error, :not_found}
    end
  end

  def list_by_organization(organization_id, tenant) do
    from(d in Document,
      join: od in OrganizationDocument,
      on: d.id == od.document_id,
      where: od.organization_id == ^organization_id and d.tenant == ^tenant,
      select: %{
        id: d.id,
        name: d.name,
        icon: d.icon,
        color: d.color,
        tenant: d.tenant,
        user_id: d.user_id,
        organization_id: od.organization_id,
        inserted_at: d.inserted_at,
        updated_at: d.updated_at
      }
    )
    |> Repo.all()
  end

  def get_document(id) do
    query =
      from(d in Document,
        join: od in OrganizationDocument,
        on: d.id == od.document_id,
        where: d.id == ^id,
        select: merge(d, %{organization_id: od.organization_id}),
        limit: 1
      )

    case Repo.one(query) do
      nil -> {:error, :not_found}
      doc -> {:ok, doc}
    end
  end

  defp initialize_y_writing(doc_id, lexical_state) do
    script_path = Application.app_dir(:realtime, "priv/scripts/convert_lexical_to_yjs")

    # Create a temporary file for the lexical state
    {:ok, temp_path} = Temp.path(%{suffix: ".json"})
    File.write!(temp_path, lexical_state)

    try do
      case System.cmd("sh", ["-c", "#{script_path} #{doc_id} @#{temp_path}"],
             stderr_to_stdout: true
           ) do
        {output, 0} ->
          Logger.info("Successfully converted lexical state to Yjs document")

          Core.Realtime.YDoc.insert_update(doc_id, output)

        {error_output, exit_code} ->
          Logger.error("Script execution failed (exit #{exit_code}): #{error_output}")
          {:error, "Could not convert lexical state to Yjs document"}
      end
    after
      File.rm(temp_path)
    end
  end

  defp maybe_insert_organization_document(multi, nil), do: multi

  defp maybe_insert_organization_document(multi, organization_id) do
    Ecto.Multi.insert(multi, :organization_document, fn %{document: document} ->
      %OrganizationDocument{}
      |> OrganizationDocument.changeset(%{
        organization_id: organization_id,
        document_id: document.id
      })
    end)
  end

  defp maybe_initialize_ydoc(multi, nil), do: multi

  defp maybe_initialize_ydoc(multi, lexical_state) do
    Ecto.Multi.run(multi, :initialize_y_writing, fn _repo, %{document: document} ->
      initialize_y_writing(document.id, lexical_state)
    end)
  end

  defp fromDto(dto) do
    %{
      "name" => name,
      "userId" => user_id,
      "tenant" => tenant,
      "body" => body,
      "icon" => icon,
      "color" => color
    } = dto

    organization_id = Map.get(dto, "organizationId", nil)

    lexical_state =
      case Map.get(dto, "lexicalState", nil) do
        nil -> nil
        value -> Jason.encode!(value)
      end

    %{
      name: name,
      body: body,
      user_id: user_id,
      tenant: tenant,
      icon: icon,
      color: color,
      lexical_state: lexical_state,
      organization_id: organization_id
    }
  end
end
