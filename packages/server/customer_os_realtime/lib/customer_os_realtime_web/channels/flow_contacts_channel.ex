defmodule CustomerOsRealtimeWeb.FlowContactsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowContacts entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "FlowContacts"
end
