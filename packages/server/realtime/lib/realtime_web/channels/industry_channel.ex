defmodule RealtimeWeb.IndustryChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Industry entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Industry"
end
