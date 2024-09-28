defmodule CustomerOsRealtimeWeb.IssuesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issues entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Issues"
end
