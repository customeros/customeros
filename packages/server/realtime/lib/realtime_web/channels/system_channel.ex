defmodule RealtimeWeb.SystemChannel do
  @moduledoc """
  This Channel broadcasts sync events to all System subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "System"
end
