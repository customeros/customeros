defmodule Web.ContactChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contact entity subscribers.
  """
  use Web.EntityChannelMacro, "Contact"
end
