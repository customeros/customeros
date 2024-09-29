defmodule RealtimeWeb.FlowContactsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowContacts entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "FlowContacts"
end
