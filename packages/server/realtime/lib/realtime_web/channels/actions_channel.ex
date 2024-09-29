defmodule RealtimeWeb.ActionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Actions entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Actions"
end
