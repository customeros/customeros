defmodule RealtimeWeb.RemindersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminders entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Reminders"
end
