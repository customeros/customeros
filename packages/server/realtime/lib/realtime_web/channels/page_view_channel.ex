defmodule RealtimeWeb.PageViewChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageView entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "PageView"
end
