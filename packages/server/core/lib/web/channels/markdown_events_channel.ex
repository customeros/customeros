defmodule Web.MarkdownEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MarkdownEvents entity subscribers.
  """
  use Web.EntityChannelMacro, "MarkdownEvents"
end
