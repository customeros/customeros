defmodule RealtimeWeb.MailBoxChannel do
  @moduledoc """
  This Channel broadcasts sync events to all MaiBoxl entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Mailbox"
end
