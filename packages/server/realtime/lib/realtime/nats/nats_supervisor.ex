defmodule Realtime.Nats.Supervisor do
  use Supervisor

  def start_link(init_args \\ []),
    do: Supervisor.start_link(__MODULE__, init_args, name: __MODULE__)

  @impl true
  def init(_init_args) do
    children = [
      Realtime.Nats.Connection,
      {Realtime.Nats.Consumer,
       stream: "webtracker",
       consumer: "realtime-webtracker-consumer",
       deliver_group: "realtime-webtracker-group",
       filter_subject: "webtracker.>",
       handler: Realtime.Nats.Handlers.WebtrackerHandler}
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
