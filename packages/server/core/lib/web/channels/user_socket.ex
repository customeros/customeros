defmodule Web.UserSocket do
  require Logger
  use Phoenix.Socket

  # A Socket handler
  #
  # It's possible to control the websocket connection and
  # assign values that can be accessed by your channel topics.

  ## Channels

  channel "organization_presence:*", Web.OrganizationViewChannel
  channel "finder:*", Web.FinderChannel
  channel "org:*", Web.TestChannel
  channel "TableViewDef:*", Web.TableViewDefChannel
  channel "TableViewDefs:*", Web.TableViewDefsChannel
  channel "Organization:*", Web.OrganizationChannel
  channel "Organizations:*", Web.OrganizationsChannel
  channel "Contract:*", Web.ContractChannel
  channel "Contracts:*", Web.ContractsChannel
  channel "Flow:*", Web.FlowChannel
  channel "Flows:*", Web.FlowsChannel
  channel "FlowSender:*", Web.FlowSenderChannel
  channel "FlowSenders:*", Web.FlowSendersChannel
  channel "ContractLineItem:*", Web.ContractLineItemChannel
  channel "ContractLineItems:*", Web.ContractLineItemsChannel
  channel "Opportunities:*", Web.OpportunitiesChannel
  channel "Opportunity:*", Web.OpportunityChannel
  channel "User:*", Web.UserChannel
  channel "Users:*", Web.UsersChannel
  channel "Invoices:*", Web.InvoicesChannel
  channel "Invoice:*", Web.InvoiceChannel
  channel "BankAccounts:*", Web.BankAccountsChannel
  channel "BankAccount:*", Web.BankAccountChannel
  channel "TenantBillingProfile:*", Web.TenantBillingProfileChannel
  channel "TenantBillingProfiles:*", Web.TenantBillingProfilesChannel
  channel "Contact:*", Web.ContactChannel
  channel "Contacts:*", Web.ContactsChannel
  channel "Reminder:*", Web.ReminderChannel
  channel "Reminders*", Web.RemindersChannel
  channel "Action:*", Web.ActionChannel
  channel "Actions:*", Web.ActionsChannel
  channel "Analysis:*", Web.AnalysisChannel
  channel "Analyses:*", Web.AnalysesChannel
  channel "InteractionEvent:*", Web.InteractionEventChannel
  channel "InteractionEvents:*", Web.InteractionEventsChannel
  channel "InteractionSession:*", Web.InteractionSessionChannel
  channel "InteractionSessions:*", Web.InteractionSessionsChannel
  channel "Issue:*", Web.IssueChannel
  channel "Issues:*", Web.IssuesChannel
  channel "LogEntry:*", Web.LogEntryChannel
  channel "LogEntries:*", Web.LogEntriesChannel
  channel "MarkdownEvent:*", Web.MarkdownEventChannel
  channel "MarkdownEvents:*", Web.MarkdownEventsChannel
  channel "Meeting:*", Web.MeetingChannel
  channel "Meetings:*", Web.MeetingsChannel
  channel "Note:*", Web.NoteChannel
  channel "Notes:*", Web.NotesChannel
  channel "Order:*", Web.OrderChannel
  channel "Orders:*", Web.OrdersChannel
  channel "PageView:*", Web.PageViewChannel
  channel "PageViews:*", Web.PageViewsChannel
  channel "Tag:*", Web.TagChannel
  channel "Tags:*", Web.TagsChannel
  channel "WorkFlow:*", Web.WorkFlowChannel
  channel "WorkFlows:*", Web.WorkFlowsChannel
  channel "FlowEmailVariables:*", Web.FlowEmailVariablesChannel
  channel "System:*", Web.SystemChannel
  channel "CustomField:*", Web.CustomFieldChannel
  channel "CustomFields:*", Web.CustomFieldsChannel
  channel "FlowParticipant:*", Web.FlowParticipantChannel
  channel "FlowParticipants:*", Web.FlowParticipantsChannel
  channel "Mailbox:*", Web.MailBoxChannel
  channel "Mailboxes:*", Web.MailboxesChannel
  channel "Agents:*", Web.AgentsChannel
  channel "Skus:*", Web.SkusChannel
  channel "Industries:*", Web.IndustriesChannel
  channel "JobRoles:*", Web.JobRolesChannel
  channel "Document:*", Web.DocumentChannel
  channel "Documents:*", Web.DocumentsChannel
  channel "Tasks:*", Web.TasksChannel
  #
  channel "OrganizationStore:*", Web.OrganizationStoreChannel

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
  #     Elixir.Web.Endpoint.broadcast("user_socket:#{user.id}", "disconnect", %{})
  #
  # Returning `nil` makes this socket anonymous.
  @impl true
  def id(_socket), do: nil
end
