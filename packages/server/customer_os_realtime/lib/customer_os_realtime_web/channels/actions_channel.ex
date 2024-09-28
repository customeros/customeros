defmodule CustomerOsRealtimeWeb.ActionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Actions entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Actions"
end
