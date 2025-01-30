defmodule RealtimeWeb.SkusChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Skus entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Skus"
end
