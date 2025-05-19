defmodule Web.PageViewsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageViews entity subscribers.
  """
  use Web.EntitiesChannelMacro, "PageViews"
end
