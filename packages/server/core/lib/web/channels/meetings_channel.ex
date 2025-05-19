defmodule Web.MeetingsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Meetings entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Meetings"
end
