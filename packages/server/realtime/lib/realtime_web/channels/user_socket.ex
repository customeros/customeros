defmodule RealtimeWeb.UserSocket do
  require Logger
  use Phoenix.Socket

  # A Socket handler
  #
  # It's possible to control the websocket connection and
  # assign values that can be accessed by your channel topics.

  ## Channels

  channel "organization_presence:*", RealtimeWeb.OrganizationViewChannel
  channel "finder:*", RealtimeWeb.FinderChannel
  channel "org:*", RealtimeWeb.TestChannel
  channel "TableViewDef:*", RealtimeWeb.TableViewDefChannel
  channel "TableViewDefs:*", RealtimeWeb.TableViewDefsChannel
  channel "Organization:*", RealtimeWeb.OrganizationChannel
  channel "Organizations:*", RealtimeWeb.OrganizationsChannel
  channel "Contract:*", RealtimeWeb.ContractChannel
  channel "Contracts:*", RealtimeWeb.ContractsChannel
  channel "Flow:*", RealtimeWeb.FlowChannel
  channel "Flows:*", RealtimeWeb.FlowsChannel
  channel "FlowContact:*", RealtimeWeb.FlowContactChannel
  channel "FlowContacts:*", RealtimeWeb.FlowContactsChannel
  channel "FlowSender:*", RealtimeWeb.FlowSenderChannel
  channel "FlowSenders:*", RealtimeWeb.FlowSendersChannel
  channel "ContractLineItem:*", RealtimeWeb.ContractLineItemChannel
  channel "ContractLineItems:*", RealtimeWeb.ContractLineItemsChannel
  channel "Opportunities:*", RealtimeWeb.OpportunitiesChannel
  channel "Opportunity:*", RealtimeWeb.OpportunityChannel
  channel "User:*", RealtimeWeb.UserChannel
  channel "Users:*", RealtimeWeb.UsersChannel
  channel "Invoices:*", RealtimeWeb.InvoicesChannel
  channel "Invoice:*", RealtimeWeb.InvoiceChannel
  channel "BankAccounts:*", RealtimeWeb.BankAccountsChannel
  channel "BankAccount:*", RealtimeWeb.BankAccountChannel
  channel "TenantBillingProfile:*", RealtimeWeb.TenantBillingProfileChannel
  channel "TenantBillingProfiles:*", RealtimeWeb.TenantBillingProfilesChannel
  channel "Contact:*", RealtimeWeb.ContactChannel
  channel "Contacts:*", RealtimeWeb.ContactsChannel
  channel "Reminder*", RealtimeWeb.ReminderChannel
  channel "Reminders*", RealtimeWeb.RemindersChannel
  channel "Action:*", RealtimeWeb.ActionChannel
  channel "Actions:*", RealtimeWeb.ActionsChannel
  channel "Analysis:*", RealtimeWeb.AnalysisChannel
  channel "Analyses:*", RealtimeWeb.AnalysesChannel
  channel "InteractionEvent:*", RealtimeWeb.InteractionEventChannel
  channel "InteractionEvents:*", RealtimeWeb.InteractionEventsChannel
  channel "InteractionSession:*", RealtimeWeb.InteractionSessionChannel
  channel "InteractionSessions:*", RealtimeWeb.InteractionSessionsChannel
  channel "Issue:*", RealtimeWeb.IssueChannel
  channel "Issues:*", RealtimeWeb.IssuesChannel
  channel "LogEntry:*", RealtimeWeb.LogEntryChannel
  channel "LogEntries:*", RealtimeWeb.LogEntriesChannel
  channel "Meeting:*", RealtimeWeb.MeetingChannel
  channel "Meetings:*", RealtimeWeb.MeetingsChannel
  channel "Note:*", RealtimeWeb.NoteChannel
  channel "Notes:*", RealtimeWeb.NotesChannel
  channel "Order:*", RealtimeWeb.OrderChannel
  channel "Orders:*", RealtimeWeb.OrdersChannel
  channel "PageView:*", RealtimeWeb.PageViewChannel
  channel "PageViews:*", RealtimeWeb.PageViewsChannel
  channel "Tag:*", RealtimeWeb.TagChannel
  channel "Tags:*", RealtimeWeb.TagsChannel
  channel "WorkFlow:*", RealtimeWeb.WorkFlowChannel
  channel "WorkFlows:*", RealtimeWeb.WorkFlowsChannel
  channel "FlowEmailVariables:*", RealtimeWeb.FlowEmailVariablesChannel
  channel "System:*", RealtimeWeb.SystemChannel

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
