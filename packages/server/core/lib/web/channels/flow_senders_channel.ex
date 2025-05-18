defmodule Web.FlowSendersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowSenders entity subscribers.
  """
  use Web.EntitiesChannelMacro, "FlowSenders"
end
