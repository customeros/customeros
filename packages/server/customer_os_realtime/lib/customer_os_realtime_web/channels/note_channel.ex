defmodule CustomerOsRealtimeWeb.NoteChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Note entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Note"
end
