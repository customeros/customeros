import React from 'react';

import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { FeaturedIcon } from '@ui/media/Icon';
import { FileX02 } from '@ui/media/icons/FileX02';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { GetInvoiceQuery } from '@shared/graphql/getInvoice.generated';
import { InvoiceSkeleton } from '@shared/components/Invoice/InvoiceSkeleton';
import {
  InvoiceLine,
  BankAccount,
  InvoiceCustomer,
  InvoiceProvider,
} from '@graphql/types';
interface InvoicePreviewModalProps {
  isError: boolean;
  isFetching: boolean;
  data: GetInvoiceQuery | undefined;
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
  isFetching,
  isError,
  data,
}) => {
  const client = getGraphQLClient();

  const { data: bankAccountsData } = useBankAccountsQuery(client);
  const { data: tenantBillingProfile } = useTenantBillingProfilesQuery(client);
  if (isFetching) {
    return <InvoiceSkeleton />;
  }

  if (!data?.invoice || isError) {
    return (
      <div className='flex flex-col items-center px-4 py-4 mt-5 overflow-hidden'>
        <FeaturedIcon colorScheme='warning'>
          <FileX02 boxSize='7' />
        </FeaturedIcon>
        <h2 className='text-md mt-4 mb-1'>Preview not available</h2>
        <span className='text-center text-sm text-gray-500'>
          Sorry, selected invoice cannot be previewed at this moment
        </span>
      </div>
    );
  }

  const customerAddressData = extractAddressData(data?.invoice?.customer);
  const providerAddressData = extractAddressData(data?.invoice?.provider);

  return (
    <Invoice
      status={data?.invoice?.status}
      invoicePeriodStart={data?.invoice?.invoicePeriodStart}
      invoicePeriodEnd={data?.invoice?.invoicePeriodEnd}
      tax={data?.invoice?.taxDue}
      note={data?.invoice?.note}
      from={providerAddressData}
      total={data?.invoice.amountDue}
      dueDate={data?.invoice.due}
      subtotal={data?.invoice.subtotal}
      issueDate={data?.invoice?.issued}
      billedTo={customerAddressData}
      invoiceNumber={data?.invoice?.invoiceNumber ?? ''}
      lines={(data?.invoice?.invoiceLineItems as Array<InvoiceLine>) ?? []}
      currency={data?.invoice?.currency || 'USD'}
      canPayWithBankTransfer={
        tenantBillingProfile?.tenantBillingProfiles?.[0]
          ?.canPayWithBankTransfer &&
        data?.invoice?.contract?.billingDetails?.canPayWithBankTransfer
      }
      availableBankAccount={
        bankAccountsData?.bankAccounts?.find(
          (e) => e.currency === data?.invoice?.currency,
        ) as BankAccount
      }
    />
  );
};
