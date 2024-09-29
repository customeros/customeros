defmodule RealtimeWeb.NotesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Notes entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Notes"
end
