defmodule CustomerOsRealtimeWeb.AnalysesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analyses entity subscribers.
  """

  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Analyses"
end
