defmodule RealtimeWeb.IndustriesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Industries entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Industries"
end
