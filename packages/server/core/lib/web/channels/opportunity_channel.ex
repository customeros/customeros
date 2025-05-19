defmodule Web.OpportunityChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunity entity subscribers.
  """
  use Web.EntityChannelMacro, "Opportunity"
end
