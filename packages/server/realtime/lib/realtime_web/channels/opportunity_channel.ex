defmodule RealtimeWeb.OpportunityChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunity entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Opportunity"
end
