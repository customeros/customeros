defmodule Web.UsersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Users entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Users"
end
