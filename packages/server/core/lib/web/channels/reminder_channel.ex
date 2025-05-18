defmodule Web.ReminderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminder entity subscribers.
  """
  use Web.EntityChannelMacro, "Reminder"
end
