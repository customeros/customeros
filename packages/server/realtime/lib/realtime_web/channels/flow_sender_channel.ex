defmodule RealtimeWeb.FlowSenderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowSender entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "FlowSender"
end
