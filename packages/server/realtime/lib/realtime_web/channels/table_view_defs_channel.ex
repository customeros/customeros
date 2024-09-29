defmodule RealtimeWeb.TableViewDefsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDefs entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "TableViewDefs"
end
