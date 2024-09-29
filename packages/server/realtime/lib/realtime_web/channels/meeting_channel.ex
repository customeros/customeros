defmodule RealtimeWeb.MeetingChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meeting entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Meeting"
end
