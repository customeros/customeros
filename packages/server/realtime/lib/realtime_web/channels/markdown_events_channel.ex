defmodule RealtimeWeb.MarkdownEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MarkdownEvents entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "MarkdownEvents"
end
