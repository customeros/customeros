defmodule Mix.Tasks.Proto.Gen do
  use Mix.Task

  @shortdoc "Generates Elixir modules from .proto files using protobuf.generate"

  def run(_args) do
    Mix.shell().info("Generating Elixir modules from .proto files...")

    output_path = "lib/proto"

    # Delete the old generated files
    Mix.shell().info("🧹 Cleaning #{output_path} ...")
    File.rm_rf!(output_path)

    # Recreate the directory
    File.mkdir_p!(output_path)

    # Collect all .proto files
    proto_files = Path.wildcard("priv/protos/**/*.proto")

    if proto_files == [] do
      Mix.raise("No .proto files found in priv/protos/")
    end

    proto_args = Enum.join(proto_files, " ")

    # Run the protobuf.generate command
    cmd = """
    mix protobuf.generate \
      --include-path priv/protos \
      --generate-descriptors=true \
      --output-path #{output_path} \
      --package-prefix=core \
      #{proto_args}
    """

    Mix.shell().cmd(cmd)
  end
end
