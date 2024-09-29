defmodule RealtimeWeb.UsersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Users entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Users"
end
