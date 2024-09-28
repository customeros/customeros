defmodule CustomerOsRealtimeWeb.OpportunityChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunity entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Opportunity"
end
