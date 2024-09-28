defmodule CustomerOsRealtimeWeb.IssueChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issue entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Issue"
end
