import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { FeaturedIcon } from '@ui/media/Icon';
import { FileX02 } from '@ui/media/icons/FileX02';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
import {
  InvoiceLine,
  BankAccount,
  InvoiceCustomer,
  InvoiceProvider,
  Invoice as TInvoice,
} from '@graphql/types';
interface InvoicePreviewModalProps {
  isFetching?: boolean;
  invoice: TInvoice | undefined | null;
}

const extractAddressData = (invoiceData: InvoiceCustomer | InvoiceProvider) => {
  return {
    zip: invoiceData?.addressZip ?? '',
    email: (invoiceData as InvoiceCustomer)?.email ?? '',
    name: invoiceData?.name ?? '',
    country: invoiceData?.addressCountry ?? '',
    locality: invoiceData?.addressLocality ?? '',
    addressLine1: invoiceData?.addressLine1 ?? '',
    addressLine2: invoiceData?.addressLine2 ?? '',
    region: invoiceData?.addressRegion ?? '',
  };
};

export const InvoicePreviewModalContent: React.FC<InvoicePreviewModalProps> = ({
  invoice,
  isFetching,
}) => {
  const client = getGraphQLClient();

  const { data: bankAccountsData } = useBankAccountsQuery(client);
  const { data: tenantBillingProfile } = useTenantBillingProfilesQuery(client);
  if (isFetching) {
    return <InvoiceSkeleton />;
  }

  if (!invoice) {
    return (
      <div className='flex flex-col items-center px-4 py-4 mt-5 overflow-hidden'>
        <FeaturedIcon colorScheme='warning'>
          <FileX02 className='size-7' />
        </FeaturedIcon>
        <h2 className='text-md mt-4 mb-1'>Preview not available</h2>
        <span className='text-center text-sm text-gray-500'>
          Sorry, selected invoice cannot be previewed at this moment
        </span>
      </div>
    );
  }

  const customerAddressData = extractAddressData(invoice?.customer);
  const providerAddressData = extractAddressData(invoice?.provider);

  return (
    <Invoice
      status={invoice?.status}
      invoicePeriodStart={invoice?.invoicePeriodStart}
      invoicePeriodEnd={invoice?.invoicePeriodEnd}
      tax={invoice?.taxDue}
      note={invoice?.note}
      from={providerAddressData}
      total={invoice.amountDue}
      dueDate={invoice.due}
      subtotal={invoice.subtotal}
      issueDate={invoice?.issued}
      billedTo={customerAddressData}
      invoiceNumber={invoice?.invoiceNumber ?? ''}
      lines={(invoice?.invoiceLineItems as Array<InvoiceLine>) ?? []}
      currency={invoice?.currency || 'USD'}
      canPayWithBankTransfer={
        tenantBillingProfile?.tenantBillingProfiles?.[0]
          ?.canPayWithBankTransfer &&
        invoice?.contract?.billingDetails?.canPayWithBankTransfer
      }
      availableBankAccount={
        bankAccountsData?.bankAccounts?.find(
          (e) => e.currency === invoice?.currency,
        ) as BankAccount
      }
    />
  );
};
