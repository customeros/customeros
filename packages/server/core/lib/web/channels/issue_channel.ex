defmodule Web.IssueChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issue entity subscribers.
  """
  use Web.EntityChannelMacro, "Issue"
end
