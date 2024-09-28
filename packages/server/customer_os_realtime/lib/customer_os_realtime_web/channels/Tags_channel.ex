defmodule CustomerOsRealtimeWeb.TagsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tags entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Tags"
end
