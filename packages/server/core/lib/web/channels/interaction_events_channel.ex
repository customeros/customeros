defmodule Web.InteractionEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvents entity subscribers.
  """
  use Web.EntitiesChannelMacro, "InteractionEvents"
end
