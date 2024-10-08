defmodule RealtimeWeb.TagChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tag entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Tag"
end
