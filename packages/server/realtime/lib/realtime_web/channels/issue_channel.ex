defmodule RealtimeWeb.IssueChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issue entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Issue"
end
