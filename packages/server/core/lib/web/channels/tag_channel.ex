defmodule Web.TagChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tag entity subscribers.
  """
  use Web.EntityChannelMacro, "Tag"
end
