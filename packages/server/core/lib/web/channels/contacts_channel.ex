defmodule Web.ContactsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contacts entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Contacts"
end
