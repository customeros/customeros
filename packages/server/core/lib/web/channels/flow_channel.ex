defmodule Web.FlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flow entity subscribers.
  """
  use Web.EntityChannelMacro, "Flow"
end
