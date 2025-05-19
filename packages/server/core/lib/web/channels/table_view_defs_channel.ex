defmodule Web.TableViewDefsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDefs entity subscribers.
  """
  use Web.EntitiesChannelMacro, "TableViewDefs"
end
