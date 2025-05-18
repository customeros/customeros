defmodule Web.CustomFieldsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all CustomFields entity subscribers.
  """
  use Web.EntitiesChannelMacro, "CustomFields"
end
