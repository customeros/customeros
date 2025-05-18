defmodule Web.TableViewDefChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TableViewDef entity subscribers.
  """
  use Web.EntityChannelMacro, "TableViewDef"
end
