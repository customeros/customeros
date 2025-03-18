defmodule Realtime.Documents do
  @moduledoc """
  The Documents context.
  """
  import Ecto.Query, warn: false
  require Logger
  alias Realtime.Repo
  alias Realtime.Documents.{Document, OrganizationDocument}

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
    |> Ecto.Multi.run(:initialize_y_writing, fn _repo, %{document: document} ->
      initialize_y_writing(document.id, payload.lexical_state)
    end)
    |> Repo.transaction()
  end

  def list_by_organization(organization_id, tenant) do
    from(d in Document,
      join: od in OrganizationDocument,
      on: d.id == od.document_id,
      where: od.organization_id == ^organization_id and d.tenant == ^tenant,
      select: %{
        id: d.id,
        name: d.name,
        tenant: d.tenant,
        user_id: d.user_id,
        inserted_at: d.inserted_at,
        updated_at: d.updated_at
      }
    )
    |> Repo.all()
  end

  defp initialize_y_writing(doc_id, lexical_state) do
    script_path = Application.app_dir(:realtime, "priv/scripts/convert_lexical_to_yjs")
    encoded_lexical_state = Jason.encode!(lexical_state)

    # Create a temporary file for the lexical state
    {:ok, temp_path} = Temp.path(%{suffix: ".json"})
    File.write!(temp_path, encoded_lexical_state)

    try do
      case System.cmd("sh", ["-c", "#{script_path} #{doc_id} @#{temp_path}"],
             stderr_to_stdout: true
           ) do
        {output, 0} ->
          Logger.info("Successfully converted lexical state to Yjs document")

          Realtime.YDoc.insert_update(doc_id, output)

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

  defp fromDto(dto) do
    %{
      "name" => name,
      "userId" => user_id,
      "tenant" => tenant,
      "body" => body,
      "lexicalState" => lexical_state
    } = dto

    organization_id = Map.get(dto, "organizationId", nil)

    %{
      name: name,
      body: body,
      user_id: user_id,
      tenant: tenant,
      lexical_state: Jason.encode!(lexical_state),
      organization_id: organization_id
    }
  end
end
