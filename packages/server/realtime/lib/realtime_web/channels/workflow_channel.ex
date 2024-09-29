defmodule RealtimeWeb.WorkFlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlow entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "WorkFlow"
end
