defmodule Realtime.Util do
  @moduledoc """
  Utility functions for Realtime.
  """

  @default_colors [
    "#F97066",
    "#F63D68",
    "#9C4221",
    "#ED8936",
    "#FDB022",
    "#ECC94B",
    "#86CB3C",
    "#38A169",
    "#3B7C0F",
    "#0BC5EA",
    "#2ED3B7",
    "#3182ce",
    "#004EEB",
    "#9E77ED",
    "#7839EE",
    "#D444F1",
    "#9F1AB1",
    "#D53F8C",
    "#98A2B3",
    "#667085",
    "#0C111D"
  ]

  def generate_random_color() do
    Enum.random(@default_colors)
  end
end
