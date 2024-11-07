import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetInvoicesQueryVariables = Types.Exact<{
  pagination: Types.Pagination;
  where?: Types.InputMaybe<Types.Filter>;
  sort?: Types.InputMaybe<Array<Types.SortBy> | Types.SortBy>;
}>;

export type GetInvoicesQuery = {
  __typename?: 'Query';
  invoices: {
    __typename?: 'InvoicesPage';
    totalElements: any;
    totalAvailable: any;
    content: Array<{
      __typename?: 'Invoice';
      issued: any;
      invoiceUrl: string;
      invoiceNumber: string;
      invoicePeriodStart: any;
      invoicePeriodEnd: any;
      due: any;
      amountDue: number;
      currency: string;
      dryRun: boolean;
      status?: Types.InvoiceStatus | null;
      preview: boolean;
      subtotal: number;
      taxDue: number;
      metadata: { __typename?: 'Metadata'; id: string; created: any };
      organization: {
        __typename?: 'Organization';
        metadata: { __typename?: 'Metadata'; id: string };
      };
      customer: {
        __typename?: 'InvoiceCustomer';
        name?: string | null;
        email?: string | null;
        addressLine1?: string | null;
        addressLine2?: string | null;
        addressZip?: string | null;
        addressLocality?: string | null;
        addressCountry?: string | null;
        addressRegion?: string | null;
      };
      contract: {
        __typename?: 'Contract';
        contractName: string;
        metadata: { __typename?: 'Metadata'; id: string };
        billingDetails?: {
          __typename?: 'BillingDetails';
          billingCycleInMonths?: any | null;
          canPayWithBankTransfer?: boolean | null;
        } | null;
      };
      provider: {
        __typename?: 'InvoiceProvider';
        logoUrl?: string | null;
        logoRepositoryFileId?: string | null;
        name?: string | null;
        addressLine1?: string | null;
        addressLine2?: string | null;
        addressZip?: string | null;
        addressLocality?: string | null;
        addressCountry?: string | null;
        addressRegion?: string | null;
      };
      invoiceLineItems: Array<{
        __typename?: 'InvoiceLine';
        quantity: any;
        subtotal: number;
        taxDue: number;
        total: number;
        price: number;
        description: string;
        metadata: {
          __typename?: 'Metadata';
          id: string;
          created: any;
          lastUpdated: any;
          source: Types.DataSource;
          sourceOfTruth: Types.DataSource;
          appSource: string;
        };
        contractLineItem: {
          __typename?: 'ServiceLineItem';
          serviceStarted: any;
          price: number;
          billingCycle: Types.BilledType;
        };
      }>;
    }>;
  };
};
