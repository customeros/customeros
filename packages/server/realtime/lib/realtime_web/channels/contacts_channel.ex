defmodule RealtimeWeb.ContactsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contacts entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Contacts"
end
