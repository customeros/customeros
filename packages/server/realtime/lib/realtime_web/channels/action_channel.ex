defmodule RealtimeWeb.ActionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Action entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Action"
end
