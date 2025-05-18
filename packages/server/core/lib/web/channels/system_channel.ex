defmodule Web.SystemChannel do
  @moduledoc """
  This Channel broadcasts sync events to all System subscribers.
  """
  use Web.EntitiesChannelMacro, "System"
end
