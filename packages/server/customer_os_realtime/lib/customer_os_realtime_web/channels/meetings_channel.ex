defmodule CustomerOsRealtimeWeb.MeetingsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meetings entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Meetings"
end
