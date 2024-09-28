defmodule CustomerOsRealtimeWeb.OpportunitiesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunities entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Opportunities"
end
