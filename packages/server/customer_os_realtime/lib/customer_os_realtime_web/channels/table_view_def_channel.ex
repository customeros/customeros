defmodule CustomerOsRealtimeWeb.TableViewDefChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDef entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "TableViewDef"
end
