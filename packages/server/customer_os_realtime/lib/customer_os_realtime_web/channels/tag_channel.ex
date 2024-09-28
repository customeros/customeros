defmodule CustomerOsRealtimeWeb.TagChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tag entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Tag"
end
