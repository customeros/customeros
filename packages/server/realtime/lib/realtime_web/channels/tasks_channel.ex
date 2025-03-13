defmodule RealtimeWeb.TasksChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tasks entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Tasks"
end
