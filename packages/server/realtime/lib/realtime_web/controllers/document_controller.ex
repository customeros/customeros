defmodule RealtimeWeb.DocumentController do
  use RealtimeWeb, :controller
  require Logger

  def create(conn, %{"docId" => doc_id, "lexicalState" => lexical_state}) do
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
          Logger.info("Script executed successfully: #{output}")

          Realtime.YDoc.insert_update(doc_id, output)

          conn
          |> put_status(:created)
          |> json(%{message: "Document created"})

        {error_output, exit_code} ->
          Logger.error("Script execution failed (exit #{exit_code}): #{error_output}")

          conn
          |> put_status(:internal_server_error)
          |> json(%{error: "Script execution failed", details: error_output})
      end
    after
      # Clean up the temporary file
      File.rm(temp_path)
    end
  end
end
