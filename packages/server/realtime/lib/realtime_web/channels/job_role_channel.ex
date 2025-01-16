defmodule RealtimeWeb.JobRoleChannel do
  @moduledoc """
  This Channel broadcasts sync events to all JobRoleChannel entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "JobRoleChannel"
end
