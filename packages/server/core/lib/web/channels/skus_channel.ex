defmodule Web.SkusChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Skus entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Skus"
end
