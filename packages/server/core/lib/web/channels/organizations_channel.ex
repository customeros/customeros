defmodule Web.OrganizationsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organizations entity subscribers.
  """

  use Web.EntitiesChannelMacro, "Organizations"
end
