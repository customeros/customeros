defmodule RealtimeWeb.JobRolesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all JobRolesChannel entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "JobRoles"
end
