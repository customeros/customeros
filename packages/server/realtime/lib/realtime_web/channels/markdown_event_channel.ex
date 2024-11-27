defmodule RealtimeWeb.MarkdownEventChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MarkdownEvent entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "MarkdownEvent"
end
