defmodule Realtime.ColorManager do
  require Logger
  use GenServer

  @moduledoc """
  Realtime.ColorManager keeps track of color assignments for users.
  It persists the color assignments in an ETS table. It also provides
  functionality to assign and release colors for users.
  """

  @name __MODULE__

  @colors [
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

  def start_link(_) do
    Logger.info("Color Manager started")
    GenServer.start_link(__MODULE__, [], name: @name)
  end

  def init(_) do
    # Initialize ETS table to store color assignments
    :ets.new(:color_assignments, [:set, :named_table])
    :ets.new(:color_assignments_counter, [:set, :named_table])
    :ets.insert(:color_assignments_counter, {:counter, 0})

    {:ok, self()}
  end

  def assign_color(user_id) do
    GenServer.call(@name, {:assign_color, user_id})
  end

  def release_color(user_id) do
    GenServer.call(@name, {:release_color, user_id})
  end

  def handle_call({:assign_color, user_id}, _from, state) do
    case :ets.lookup(:color_assignments, user_id) do
      [] ->
        color = get_next_color()
        :ets.insert(:color_assignments, {user_id, color})
        Logger.info("Assigned color #{color} to user #{user_id}")
        {:reply, {:ok, color}, state}

      [{^user_id, color}] ->
        {:reply, {:ok, color}, state}
    end
  end

  def handle_call({:release_color, user_id}, _from, state) do
    :ets.delete(:color_assignments, user_id)
    {:reply, :ok, state}
  end

  defp get_next_color() do
    :ets.update_counter(:color_assignments_counter, :counter, 1)
    {_, count} = hd(:ets.lookup(:color_assignments_counter, :counter))

    Enum.at(@colors, rem(count, length(@colors)))
  end
end
