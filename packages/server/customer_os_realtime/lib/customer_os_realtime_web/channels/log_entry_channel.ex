defmodule CustomerOsRealtimeWeb.LogEntryChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntry entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "LogEntry"
end
