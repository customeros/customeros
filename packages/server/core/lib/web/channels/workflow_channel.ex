defmodule Web.WorkFlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlow entity subscribers.
  """
  use Web.EntityChannelMacro, "WorkFlow"
end
