defmodule CustomerOsRealtimeWeb.MeetingChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meeting entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Meeting"
end
