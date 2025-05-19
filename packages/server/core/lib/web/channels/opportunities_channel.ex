defmodule Web.OpportunitiesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Opportunities entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Opportunities"
end
