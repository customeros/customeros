defmodule RealtimeWeb.MailboxesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Mailboxes entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Mailboxes"
end
