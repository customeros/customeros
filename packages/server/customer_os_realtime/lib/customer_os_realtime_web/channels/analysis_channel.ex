defmodule CustomerOsRealtimeWeb.AnalysisChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analysis entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Analysis"
end
