defmodule Web.InteractionEventChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvent entity subscribers.
  """
  use Web.EntityChannelMacro, "InteractionEvent"
end
