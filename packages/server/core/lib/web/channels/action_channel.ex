defmodule Web.ActionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Action entity subscribers.
  """
  use Web.EntityChannelMacro, "Action"
end
