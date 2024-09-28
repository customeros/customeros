defmodule CustomerOsRealtimeWeb.ReminderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminder entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Reminder"
end
