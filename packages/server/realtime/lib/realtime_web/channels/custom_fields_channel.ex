defmodule RealtimeWeb.CustomFieldsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all CustomFields entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "CustomFields"
end
