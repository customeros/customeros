defmodule CustomerOsRealtimeWeb.TableViewDefsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDefs entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "TableViewDefs"
end
