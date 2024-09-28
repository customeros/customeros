defmodule CustomerOsRealtimeWeb.PageViewsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageViews entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "PageViews"
end
