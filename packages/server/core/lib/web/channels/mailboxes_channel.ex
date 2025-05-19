defmodule Web.MailboxesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Mailboxes entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Mailboxes"
end
