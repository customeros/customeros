defmodule RealtimeWeb.FlowSendersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowSenders entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "FlowSenders"
end
