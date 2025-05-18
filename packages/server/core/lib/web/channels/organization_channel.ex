defmodule Web.OrganizationChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organization entity subscribers.
  """
  use Web.EntityChannelMacro, "Organization"
end
