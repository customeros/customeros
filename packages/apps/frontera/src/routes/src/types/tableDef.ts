import { Node, Scalars } from '@graphql/types';
export enum ColumnViewType {
  InvoicesAmount = 'INVOICES_AMOUNT',
  InvoicesBillingCycle = 'INVOICES_BILLING_CYCLE',
  InvoicesContract = 'INVOICES_CONTRACT',
  InvoicesDueDate = 'INVOICES_DUE_DATE',
  InvoicesInvoiceNumber = 'INVOICES_INVOICE_NUMBER',
  InvoicesInvoicePreview = 'INVOICES_INVOICE_PREVIEW',
  InvoicesInvoiceStatus = 'INVOICES_INVOICE_STATUS',
  InvoicesIssueDate = 'INVOICES_ISSUE_DATE',
  InvoicesIssueDatePast = 'INVOICES_ISSUE_DATE_PAST',
  InvoicesPaymentStatus = 'INVOICES_PAYMENT_STATUS',
  OrganizationsAvatar = 'ORGANIZATIONS_AVATAR',
  OrganizationsForecastArr = 'ORGANIZATIONS_FORECAST_ARR',
  OrganizationsLastTouchpoint = 'ORGANIZATIONS_LAST_TOUCHPOINT',
  OrganizationsName = 'ORGANIZATIONS_NAME',
  OrganizationsOnboardingStatus = 'ORGANIZATIONS_ONBOARDING_STATUS',
  OrganizationsOwner = 'ORGANIZATIONS_OWNER',
  OrganizationsRelationship = 'ORGANIZATIONS_RELATIONSHIP',
  OrganizationsRenewalLikelihood = 'ORGANIZATIONS_RENEWAL_LIKELIHOOD',
  OrganizationsRenewlDate = 'ORGANIZATIONS_RENEWL_DATE',
  OrganizationsWebsite = 'ORGANIZATIONS_WEBSITE',
  RenewalsAvatar = 'RENEWALS_AVATAR',
  RenewalsForecastArr = 'RENEWALS_FORECAST_ARR',
  RenewalsLastTouchpoint = 'RENEWALS_LAST_TOUCHPOINT',
  RenewalsName = 'RENEWALS_NAME',
  RenewalsOwner = 'RENEWALS_OWNER',
  RenewalsRenewalDate = 'RENEWALS_RENEWAL_DATE',
  RenewalsRenewalLikelihood = 'RENEWALS_RENEWAL_LIKELIHOOD',
}

export type TableViewDef = Node & {
  tableType: TableViewType;
  columns: Array<ColumnView>;
  __typename?: 'TableViewDef';
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  icon: Scalars['String']['output'];
  name: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  filters: Scalars['String']['output'];
  sorting: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type ColumnView = {
  __typename?: 'ColumnView';
  columnType: ColumnViewType;
  width: Scalars['Int']['output'];
};

export enum TableViewType {
  Invoices = 'INVOICES',
  Organizations = 'ORGANIZATIONS',
  Renewals = 'RENEWALS',
}
