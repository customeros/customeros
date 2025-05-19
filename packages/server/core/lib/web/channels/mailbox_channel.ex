defmodule Web.MailBoxChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MaiBoxl entity subscribers.
  """
  use Web.EntityChannelMacro, "Mailbox"
end
