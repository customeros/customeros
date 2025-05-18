defmodule Web.OrganizationStoreChannel do
  @moduledoc """
  This channel broadcasts sync events to entity subscribers.
  It is used to dynamically handle multiple entity types.
  """
  require Logger
  use Web, :channel

  @impl true
  def join(
        _topic,
        %{"user_id" => user_id, "username" => username},
        socket
      ) do
    socket =
      socket
      |> assign(:user_id, user_id)
      |> assign(:username, username)

    send(self(), :after_join)
    {:ok, %{}, socket}
  end

  @impl true
  def join(_, _, _) do
    {:ok, self()}
  end

  @impl true
  def handle_info(:after_join, socket) do
    {:noreply, socket}
  end

  @impl true
  def handle_in("store:set", payload, socket) do
    broadcast!(socket, "store:set", payload)
    {:reply, {:ok, %{}}, socket}
  end

  @impl true
  def handle_in("store:delete", payload, socket) do
    broadcast!(socket, "store:delete", payload)
    {:reply, {:ok, %{}}, socket}
  end

  @impl true
  def handle_in("store:clear", payload, socket) do
    broadcast!(socket, "store:clear", payload)
    {:reply, {:ok, %{}}, socket}
  end

  @impl true
  def handle_in("store:invalidate", payload, socket) do
    broadcast!(socket, "store:invalidate", payload)
    {:reply, {:ok, %{}}, socket}
  end

  @impl true
  def terminate(_, socket) do
    {:ok, socket}
  end
end
