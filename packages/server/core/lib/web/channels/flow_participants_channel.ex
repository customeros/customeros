defmodule Web.FlowParticipantsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowParticipants entity subscribers.
  """
  use Web.EntitiesChannelMacro, "FlowParticipants"
end
