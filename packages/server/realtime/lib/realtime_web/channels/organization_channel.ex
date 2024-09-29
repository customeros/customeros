defmodule RealtimeWeb.OrganizationChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organization entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Organization"
end
