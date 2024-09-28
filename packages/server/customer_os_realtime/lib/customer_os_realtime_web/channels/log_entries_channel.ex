defmodule CustomerOsRealtimeWeb.LogEntriesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntries entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "LogEntries"
end
