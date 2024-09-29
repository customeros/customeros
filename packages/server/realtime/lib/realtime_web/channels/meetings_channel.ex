defmodule RealtimeWeb.MeetingsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meetings entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Meetings"
end
