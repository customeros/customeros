# defmodule Nats.Supervisor do
#   use Supervisor
#   require Logger
#
#   def start_link(init_args \\ []),
#     do: Supervisor.start_link(__MODULE__, init_args, name: __MODULE__)
#
#   @impl true
#   def init(_init_args) do
#     children = [
#       {Nats.Consumer,
#        stream: "ai",
#        consumer: "ai-consumer",
#        deliver_group: "ai-group",
#        filter_subject: "ai.>",
#        handler: Nats.Handlers.WebtrackerHandler}
#     ]
#
#     Supervisor.init(children, strategy: :one_for_one)
#   end
# end
