defmodule RealtimeWeb.FlowParticipantsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowParticipants entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "FlowParticipants"
end
