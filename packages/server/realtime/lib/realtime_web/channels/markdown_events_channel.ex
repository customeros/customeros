defmodule RealtimeWeb.MarkdownEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntry entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "MarkdownEvents"
end
