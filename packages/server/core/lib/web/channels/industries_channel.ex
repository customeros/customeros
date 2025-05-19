defmodule Web.IndustriesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Industries entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Industries"
end
