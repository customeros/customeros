defmodule Web.JobRolesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all JobRoles entity subscribers.
  """
  use Web.EntitiesChannelMacro, "JobRoles"
end
