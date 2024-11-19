defmodule RealtimeWeb.FlowParticipantChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowParticipant entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "FlowParticipant"
end
