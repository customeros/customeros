defmodule Web.CustomFieldChannel do
  @moduledoc """
  This Channel broadcasts sync events to all CustomField entity subscribers.
  """
  use Web.EntityChannelMacro, "CustomField"
end
