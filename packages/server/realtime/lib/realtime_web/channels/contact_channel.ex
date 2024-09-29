defmodule RealtimeWeb.ContactChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contact entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Contact"
end
