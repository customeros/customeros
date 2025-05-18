defmodule Web.NoteChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Note entity subscribers.
  """
  use Web.EntityChannelMacro, "Note"
end
