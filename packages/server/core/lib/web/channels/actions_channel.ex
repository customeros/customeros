defmodule Web.ActionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Actions entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Actions"
end
