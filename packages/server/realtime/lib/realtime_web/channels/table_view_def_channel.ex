defmodule RealtimeWeb.TableViewDefChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDef entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "TableViewDef"
end
