defmodule RealtimeWeb.PageViewsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageViews entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "PageViews"
end
