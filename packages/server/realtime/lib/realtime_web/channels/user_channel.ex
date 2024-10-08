defmodule RealtimeWeb.UserChannel do
  @moduledoc """
  This Channel broadcasts sync events to all User entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "User"
end
