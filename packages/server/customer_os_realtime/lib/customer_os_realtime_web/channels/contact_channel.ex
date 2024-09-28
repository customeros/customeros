defmodule CustomerOsRealtimeWeb.ContactChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contact entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Contact"
end
