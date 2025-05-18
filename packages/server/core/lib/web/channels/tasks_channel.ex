defmodule Web.TasksChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tasks entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Tasks"
end
