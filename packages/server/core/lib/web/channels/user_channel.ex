defmodule Web.UserChannel do
  @moduledoc """
  This Channel broadcasts sync events to all User entity subscribers.
  """
  use Web.EntityChannelMacro, "User"
end
