defmodule CustomerOsRealtimeWeb.EntitiesChannelMacro do
  defmacro __using__(entity_prefix) do
    quote do
      use Phoenix.Channel

      alias CustomerOsRealtimeWeb.GenericMultiChannel

      @impl true
      def join(unquote(entity_prefix) <> ":" <> entity_id, params, socket) do
        GenericMultiChannel.handle_join(unquote(entity_prefix), entity_id, params, socket)
      end

      @impl true
      def handle_info(message, socket) do
        GenericMultiChannel.handle_info(message, socket)
      end

      @impl true
      def handle_in(event, payload, socket) do
        GenericMultiChannel.handle_in(event, payload, socket)
      end

      @impl true
      def terminate(reason, socket) do
        GenericMultiChannel.terminate(reason, socket)
      end
    end
  end
end
