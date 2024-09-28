defmodule CustomerOsRealtimeWeb.PageViewChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageView entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "PageView"
end
