defmodule Web.FlowParticipantChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowParticipant entity subscribers.
  """
  use Web.EntityChannelMacro, "FlowParticipant"
end
