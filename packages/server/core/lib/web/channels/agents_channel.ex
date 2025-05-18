defmodule Web.AgentsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Agents entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Agents"
end
