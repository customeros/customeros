defmodule CustomerOsRealtimeWeb.NotesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Notes entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Notes"
end
