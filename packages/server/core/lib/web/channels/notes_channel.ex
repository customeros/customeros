defmodule Web.NotesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Notes entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Notes"
end
