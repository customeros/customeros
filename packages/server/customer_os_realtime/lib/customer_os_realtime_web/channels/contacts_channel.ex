defmodule CustomerOsRealtimeWeb.ContactsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contacts entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Contacts"
end
