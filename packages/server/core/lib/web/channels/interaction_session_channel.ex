defmodule Web.InteractionSessionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSession entity subscribers.
  """
  use Web.EntityChannelMacro, "InteractionSession"
end
