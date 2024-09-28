defmodule CustomerOsRealtimeWeb.UserChannel do
  @moduledoc """
  This Channel broadcasts sync events to all User entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "User"
end
