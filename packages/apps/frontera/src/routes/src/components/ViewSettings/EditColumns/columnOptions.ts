import { ColumnViewType } from '@graphql/types';

type InvoicesColumnType =
  | ColumnViewType.InvoicesAmount
  | ColumnViewType.InvoicesBillingCycle
  | ColumnViewType.InvoicesContract
  | ColumnViewType.InvoicesDueDate
  | ColumnViewType.InvoicesIssueDatePast
  | ColumnViewType.InvoicesInvoicePreview
  | ColumnViewType.InvoicesIssueDate
  | ColumnViewType.InvoicesInvoiceStatus
  | ColumnViewType.InvoicesPaymentStatus
  | ColumnViewType.InvoicesInvoiceNumber;

export const invoicesOptionsMap: Record<InvoicesColumnType | string, string> = {
  [ColumnViewType.InvoicesAmount]: 'Amount',
  [ColumnViewType.InvoicesBillingCycle]: 'Billing cycle',
  [ColumnViewType.InvoicesContract]: 'Contract',
  [ColumnViewType.InvoicesDueDate]: 'Due date',
  [ColumnViewType.InvoicesInvoicePreview]: 'Invoice preview',
  [ColumnViewType.InvoicesInvoiceNumber]: 'Invoice',
  [ColumnViewType.InvoicesIssueDate]: 'Issue date',
  [ColumnViewType.InvoicesIssueDatePast]: 'Issue date',
  [ColumnViewType.InvoicesInvoiceStatus]: 'Invoice status',
  [ColumnViewType.InvoicesPaymentStatus]: 'Payment status',
};

export const invoicesHelperTextMap: Record<
  InvoicesColumnType | string,
  string
> = {
  [ColumnViewType.InvoicesAmount]: 'E.g. $6,450',
  [ColumnViewType.InvoicesBillingCycle]: 'E.g. Monthly',
  [ColumnViewType.InvoicesContract]: 'E.g. Pile Contract',
  [ColumnViewType.InvoicesDueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesInvoicePreview]: 'E.g. RKD-04025',
  [ColumnViewType.InvoicesInvoiceNumber]: 'E.g. RKD-04025',
  [ColumnViewType.InvoicesIssueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesIssueDatePast]: 'E.g. 15 Jun 2019',
  [ColumnViewType.InvoicesInvoiceStatus]: 'E.g. Scheduled',
  [ColumnViewType.InvoicesPaymentStatus]: 'E.g. Paid',
};

type RenewalsColumnType =
  | ColumnViewType.RenewalsAvatar
  | ColumnViewType.RenewalsForecastArr
  | ColumnViewType.RenewalsLastTouchpoint
  | ColumnViewType.RenewalsName
  | ColumnViewType.RenewalsOwner
  | ColumnViewType.RenewalsRenewalDate
  | ColumnViewType.RenewalsRenewalLikelihood;

export const renewalsOptionsMap: Record<RenewalsColumnType | string, string> = {
  [ColumnViewType.RenewalsAvatar]: 'Avatar',
  [ColumnViewType.RenewalsForecastArr]: 'ARR Forecast',
  [ColumnViewType.RenewalsLastTouchpoint]: 'Last touchpoint',
  [ColumnViewType.RenewalsName]: 'Name',
  [ColumnViewType.RenewalsOwner]: 'Owner',
  [ColumnViewType.RenewalsRenewalDate]: 'Next Renewal',
  [ColumnViewType.RenewalsRenewalLikelihood]: 'Health',
};

export const renewalsHelperTextMap: Record<
  RenewalsColumnType | string,
  string
> = {
  [ColumnViewType.RenewalsAvatar]: 'E.g. Avatar',
  [ColumnViewType.RenewalsForecastArr]: 'E.g. $6,450',
  [ColumnViewType.RenewalsLastTouchpoint]: 'E.g. Issue updated',
  [ColumnViewType.RenewalsName]: 'E.g. Pile Contract',
  [ColumnViewType.RenewalsOwner]: 'E.g. Howard Hu',
  [ColumnViewType.RenewalsRenewalDate]: 'E.g. 1 month',
  [ColumnViewType.RenewalsRenewalLikelihood]: 'E.g. High',
};
