defmodule RealtimeWeb.IssuesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issues entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Issues"
end
