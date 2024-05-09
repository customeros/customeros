import {
  TableIdType,
  TableViewDef,
  TableViewType,
  ColumnViewType,
} from '@graphql/types';

export const mockedTableDefs: TableViewDef[] = [
  {
    id: '1',
    tableId: TableIdType.Customers,
    order: 0,
    name: 'Monthly renewals',
    icon: '',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
    tableType: TableViewType.Renewals,
    filters: '',
    sorting: '',
    columns: [
      {
        columnType: ColumnViewType.RenewalsAvatar,
        width: 32,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsName,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsForecastArr,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalLikelihood,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsOwner,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsLastTouchpoint,
        width: 100,
        visible: true,
      },
    ],
  },
  {
    id: '2',
    tableId: TableIdType.QuarterlyRenewals,

    order: 1,
    name: 'Quarterly renewals',
    icon: '',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
    tableType: TableViewType.Renewals,
    filters: '',
    sorting: '',
    columns: [
      {
        columnType: ColumnViewType.RenewalsAvatar,
        width: 32,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsName,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsForecastArr,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalLikelihood,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsOwner,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsLastTouchpoint,
        width: 100,
        visible: true,
      },
    ],
  },
  {
    id: '3',
    order: 2,
    name: 'Annual renewals',
    icon: '',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
    tableType: TableViewType.Renewals,
    tableId: TableIdType.AnnualRenewals,
    filters: '',
    sorting: '',
    columns: [
      {
        columnType: ColumnViewType.RenewalsAvatar,
        width: 32,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsName,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsForecastArr,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsRenewalLikelihood,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsOwner,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.RenewalsLastTouchpoint,
        width: 100,
        visible: true,
      },
    ],
  },
  {
    id: '4',
    order: 3,
    name: 'Upcoming invoices',
    icon: '',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
    tableType: TableViewType.Invoices,
    tableId: TableIdType.UpcomingInvoices,
    filters: '',
    sorting: '',
    columns: [
      {
        columnType: ColumnViewType.InvoicesInvoicePreview,
        width: 32,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesContract,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesBillingCycle,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesIssueDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesDueDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesAmount,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesInvoiceStatus,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesIssueDatePast,
        width: 100,
        visible: false,
      },
      {
        columnType: ColumnViewType.InvoicesPaymentStatus,
        width: 100,
        visible: false,
      },
    ],
  },
  {
    id: '5',
    order: 4,
    name: 'Issued invoices',
    icon: '',
    createdAt: '2021-08-10T14:00:00.000Z',
    updatedAt: '2021-08-10T14:00:00.000Z',
    tableType: TableViewType.Invoices,
    tableId: TableIdType.PastInvoices,
    filters: '',
    sorting: '',
    columns: [
      {
        columnType: ColumnViewType.InvoicesInvoicePreview,
        width: 32,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesContract,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesBillingCycle,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesIssueDatePast,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesDueDate,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesAmount,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesPaymentStatus,
        width: 100,
        visible: true,
      },
      {
        columnType: ColumnViewType.InvoicesIssueDate,
        width: 100,
        visible: false,
      },
      {
        columnType: ColumnViewType.InvoicesInvoiceStatus,
        width: 100,
        visible: false,
      },
    ],
  },
];
