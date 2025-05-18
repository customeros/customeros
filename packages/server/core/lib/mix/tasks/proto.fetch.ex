defmodule Mix.Tasks.Proto.Fetch do
  use Mix.Task

  @shortdoc "Fetches .proto files from GitHub"

  def run(_args) do
    Mix.shell().info("Fetching .proto files from GitHub using curl...")

    output_path = "priv/protos"

    Mix.shell().info("🧹 Cleaning #{output_path} ...")
    File.rm_rf!(output_path)

    File.mkdir_p!(output_path)

    {json_body, 0} =
      System.cmd("curl", [
        "-sSL",
        "-H",
        "Accept: application/vnd.github.v3+json",
        "https://api.github.com/repos/customeros/customeros/contents/packages/server/proto"
      ])

    files = Jason.decode!(json_body)

    for %{"name" => name, "download_url" => download_url} <- files,
        String.ends_with?(name, ".proto") do
      dest = "#{output_path}/#{name}"
      Mix.shell().info("Downloading #{name}...")

      {_out, 0} = System.cmd("curl", ["-sSL", download_url, "-o", dest])
    end
  end
end
