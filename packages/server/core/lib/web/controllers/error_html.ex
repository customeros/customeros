defmodule Web.ErrorHTML do
  use Web, :html

  embed_templates "error_html/*"
end
