defmodule RealtimeWeb.TagsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tags entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Tags"
end
