defmodule Web.MeetingChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meeting entity subscribers.
  """
  use Web.EntityChannelMacro, "Meeting"
end
