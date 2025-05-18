defmodule Web.LogEntriesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all LogEntries entity subscribers.
  """
  use Web.EntitiesChannelMacro, "LogEntries"
end
