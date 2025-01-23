defmodule RealtimeWeb.JobRolesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all JobRoles entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "JobRoles"
end
