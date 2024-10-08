defmodule RealtimeWeb.InteractionSessionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSession entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "InteractionSession"
end
