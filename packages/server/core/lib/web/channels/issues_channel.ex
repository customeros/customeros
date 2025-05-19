defmodule Web.IssuesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Issues entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Issues"
end
