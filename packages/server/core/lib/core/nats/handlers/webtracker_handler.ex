defmodule Core.Nats.Handlers.WebtrackerHandler do
  @moduledoc """
    This module handles decoding the protobuf webtracker stream messages.
  """

  alias Core.Realtime.Pb.{WebtrackerVisitorIdentified}

  def handle_message(%{body: body, topic: "webtracker.visitor.identified"}) do
    try do
      _event = WebtrackerVisitorIdentified.decode(body)
      :ok
    rescue
      e ->
        IO.warn("Decode failed: #{Exception.message(e)}")
        :error
    end
  end

  def handle_message(_) do
    {:ok}
  end
end
