defmodule CustomerOsRealtimeWeb.EntityChannelMacro do
  defmacro __using__(entity_prefix) do
    quote do
      use Phoenix.Channel

      alias CustomerOsRealtimeWeb.GenericChannel

      @impl true
      def join(unquote(entity_prefix) <> ":" <> entity_id, params, socket) do
        GenericChannel.handle_join(unquote(entity_prefix), entity_id, params, socket)
      end

      @impl true
      def handle_info(message, socket) do
        GenericChannel.handle_info(message, socket)
      end

      @impl true
      def handle_in(event, payload, socket) do
        GenericChannel.handle_in(event, payload, socket)
      end

      @impl true
      def terminate(reason, socket) do
        GenericChannel.terminate(reason, socket)
      end
    end
  end
end
