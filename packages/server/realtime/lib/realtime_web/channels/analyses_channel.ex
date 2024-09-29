defmodule RealtimeWeb.AnalysesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analyses entity subscribers.
  """

  use RealtimeWeb.EntitiesChannelMacro, "Analyses"
end
