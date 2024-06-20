[
  import_deps: [:phoenix],
  plugins: [Phoenix.LiveView.HTMLFormatter],
  locals_without_parens: [from: 2],
  inputs: ["mix.exs", "*.{heex,ex,exs}", "{config,lib,test}/**/*.{heex,ex,exs}"]
]
