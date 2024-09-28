defmodule CustomerOsRealtimeWeb.RemindersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminders entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Reminders"
end
