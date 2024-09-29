defmodule RealtimeWeb.LogEntriesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntries entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "LogEntries"
end
