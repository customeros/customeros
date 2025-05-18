defmodule Web.AnalysesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Analyses entity subscribers.
  """

  use Web.EntitiesChannelMacro, "Analyses"
end
