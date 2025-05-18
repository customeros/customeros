defmodule Web.LogEntryChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntry entity subscribers.
  """
  use Web.EntityChannelMacro, "LogEntry"
end
