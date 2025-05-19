defmodule Web.AnalysisChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analysis entity subscribers.
  """
  use Web.EntityChannelMacro, "Analysis"
end
