defmodule RealtimeWeb.NoteChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Note entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Note"
end
