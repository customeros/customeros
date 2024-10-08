defmodule RealtimeWeb.OpportunitiesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunities entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Opportunities"
end
