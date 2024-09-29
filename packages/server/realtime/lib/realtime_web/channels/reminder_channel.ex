defmodule RealtimeWeb.ReminderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminder entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Reminder"
end
