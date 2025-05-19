defmodule Web.MarkdownEventChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MarkdownEvent entity subscribers.
  """
  use Web.EntityChannelMacro, "MarkdownEvent"
end
