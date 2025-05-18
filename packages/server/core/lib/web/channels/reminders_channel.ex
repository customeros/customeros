defmodule Web.RemindersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Reminders entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Reminders"
end
