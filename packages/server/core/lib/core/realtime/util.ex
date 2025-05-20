defmodule Core.Realtime.Util do
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

  def to_camel_case_map(map) when is_map(map) do
    map
    |> Map.delete(:__meta__)
    |> Enum.map(fn {key, value} ->
      new_key =
        key |> to_string() |> Macro.camelize() |> String.replace(~r/^./, &String.downcase/1)

      {new_key, transform_value(value)}
    end)
    |> Enum.into(%{})
  end

  defp transform_value(%DateTime{} = datetime), do: datetime
  defp transform_value(value) when is_map(value), do: to_camel_case_map(value)
  defp transform_value(value) when is_list(value), do: Enum.map(value, &transform_value/1)
  defp transform_value(value), do: value
end
