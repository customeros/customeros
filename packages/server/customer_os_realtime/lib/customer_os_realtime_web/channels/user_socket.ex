defmodule CustomerOsRealtimeWeb.UserSocket do
  require Logger
  use Phoenix.Socket

  # A Socket handler
  #
  # It's possible to control the websocket connection and
  # assign values that can be accessed by your channel topics.

  ## Channels

  channel "organization_presence:*", CustomerOsRealtimeWeb.OrganizationViewChannel
  channel "finder:*", CustomerOsRealtimeWeb.FinderChannel
  channel "org:*", CustomerOsRealtimeWeb.TestChannel
  channel "TableViewDef:*", CustomerOsRealtimeWeb.TableViewDefChannel
  channel "TableViewDefs:*", CustomerOsRealtimeWeb.TableViewDefsChannel
  channel "Organization:*", CustomerOsRealtimeWeb.OrganizationChannel
  channel "Organizations:*", CustomerOsRealtimeWeb.OrganizationsChannel
  channel "Contract:*", CustomerOsRealtimeWeb.ContractChannel
  channel "Contracts:*", CustomerOsRealtimeWeb.ContractsChannel
  channel "ContractLineItem:*", CustomerOsRealtimeWeb.ContractLineItemChannel
  channel "ContractLineItems:*", CustomerOsRealtimeWeb.ContractLineItemsChannel
  channel "Opportunities:*", CustomerOsRealtimeWeb.OpportunitiesChannel
  channel "Opportunity:*", CustomerOsRealtimeWeb.OpportunityChannel
  channel "User:*", CustomerOsRealtimeWeb.UserChannel
  channel "Users:*", CustomerOsRealtimeWeb.UsersChannel
  channel "Invoices:*", CustomerOsRealtimeWeb.InvoicesChannel
  channel "Invoice:*", CustomerOsRealtimeWeb.InvoiceChannel

  # Socket params are passed from the client and can
  # be used to verify and authenticate a user. After
  # verification, you can put default assigns into
  # the socket that will be set for all channels, ie
  #
  #     {:ok, assign(socket, :user_id, verified_user_id)}
  #
  # To deny connection, return `:error` or `{:error, term}`. To control the
  # response the client receives in that case, [define an error handler in the
  # websocket
  # configuration](https://hexdocs.pm/phoenix/Phoenix.Endpoint.html#socket/3-websocket-configuration).
  #
  # See `Phoenix.Token` documentation for examples in
  # performing token verification on connect.

  # @impl true
  # def connect(_params, socket, _connect_info) do
  #   Logger.info "Reached connect in user_socket.ex"
  #   {:ok, socket}
  # end

  @impl true
  def connect(params, socket, _connect_info) do
    if authorized(params["token"]) do
      {:ok, socket}
    else
      {:error, %{reason: "unauthorized"}}
    end
  end

  defp authorized(token) do
    token == System.get_env("API_TOKEN")
  end

  # Socket id's are topics that allow you to identify all sockets for a given user:
  #
  #     def id(socket), do: "user_socket:#{socket.assigns.user_id}"
  #
  # Would allow you to broadcast a "disconnect" event and terminate
  # all active sockets and channels for a given user:
  #
  #     Elixir.RealtimeWeb.Endpoint.broadcast("user_socket:#{user.id}", "disconnect", %{})
  #
  # Returning `nil` makes this socket anonymous.
  @impl true
  def id(_socket), do: nil
end
