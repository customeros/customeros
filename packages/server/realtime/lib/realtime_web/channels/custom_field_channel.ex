defmodule RealtimeWeb.CustomFieldChannel do
  @moduledoc """
  This Channel broadcasts sync events to all CustomField entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "CustomField"
end
