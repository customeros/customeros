defmodule RealtimeWeb.AnalysisChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analysis entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Analysis"
end
