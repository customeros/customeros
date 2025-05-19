defmodule Web.TagsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Tags entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Tags"
end
