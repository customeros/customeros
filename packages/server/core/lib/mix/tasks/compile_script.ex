defmodule Mix.Tasks.Compile.Script do
  use Mix.Task

  @shortdoc "Compiles the TypeScript script using Bun"

  def run(_) do
    if running_in_elixirls?() do
      Mix.shell().info("Skipping TypeScript compilation in ElixirLS")
    end

    Mix.shell().info("Compiling TypeScript script with Bun...")

    # Define script and output paths
    script_dir = Path.expand("scripts")
    output_dir = Path.expand("priv/scripts")

    script_path = Path.join(script_dir, "convert_lexical_to_yjs.ts")
    output_path = Path.join(output_dir, "convert_lexical_to_yjs")

    # Ensure priv/scripts exists
    File.mkdir_p!(output_dir)

    # Run the Bun build command and capture output
    command = "cd #{script_dir} && bun build #{script_path} --compile --outfile=#{output_path}"

    {output, exit_code} = System.cmd("sh", ["-c", command], stderr_to_stdout: true)

    # Print output for debugging
    Mix.shell().info(output)

    if exit_code == 0 do
      Mix.shell().info("✅ Successfully compiled script to #{output_path}")
    else
      Mix.raise("❌ Failed to compile TS script.\n\nError output:\n#{output}")
    end
  end

  defp running_in_elixirls? do
    System.get_env("ELIXIR_LS") == "true"
  end
end
