defmodule Web.FlowSenderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowSender entity subscribers.
  """
  use Web.EntityChannelMacro, "FlowSender"
end
