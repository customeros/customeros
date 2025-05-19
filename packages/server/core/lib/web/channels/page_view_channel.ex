defmodule Web.PageViewChannel do
  @moduledoc """
  This Channel broadcasts sync events to all PageView entity subscribers.
  """
  use Web.EntityChannelMacro, "PageView"
end
