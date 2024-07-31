import { useMemo } from 'react';

import { TenantBillingDetailsDto } from '@settings/components/Tabs/panels/BillingPanel/TenantBillingProfile.dto';
import { useBankTransferSelectionContext } from '@settings/components/Tabs/panels/BillingPanel/context/BankTransferSelectionContext';

import { Invoice } from '@shared/components/Invoice/Invoice';
import { DataSource, InvoiceLine, InvoiceStatus } from '@graphql/types';

export const BillingPanelInvoice = ({
  values,
  isInvoiceProviderFocused,
  isInvoiceProviderDetailsHovered,
}: {
  values?: TenantBillingDetailsDto;
  isInvoiceProviderFocused?: boolean;
  isInvoiceProviderDetailsHovered?: boolean;
}) => {
  const { defaultSelectedAccount, hoveredAccount, focusedAccount } =
    useBankTransferSelectionContext();

  const invoicePreviewStaticData = useMemo(
    () => ({
      status: InvoiceStatus.Scheduled,
      invoiceNumber: 'INV-003',
      lines: [
        {
          subtotal: 100,
          createdAt: new Date().toISOString(),
          metadata: {
            id: 'dummy-id',
            created: new Date().toISOString(),
            lastUpdated: new Date().toISOString(),
            source: DataSource.Openline,
            sourceOfTruth: DataSource.Openline,
            appSource: DataSource.Openline,
          },
          description: 'Professional tier',
          price: 50,
          quantity: 2,
          total: 100,
          taxDue: 0,
        } as unknown as InvoiceLine,
      ],
      tax: 0,
      note: '',
      total: 100,
      dueDate: new Date().toISOString(),
      subtotal: 100,
      issueDate: new Date().toISOString(),
      billedTo: {
        addressLine1: '29 Maple Lane',
        addressLine2: 'Springfield, Haven County',
        locality: 'San Francisco',
        region: 'CA',
        zip: '89302',
        country: 'United States of America',
        email: 'invoices@acme.com',
        name: 'Acme Corp.',
      },
    }),
    [],
  );

  const displayedBankAccount =
    focusedAccount || hoveredAccount || defaultSelectedAccount;

  return (
    <div className='border-r border-gray-300 max-h-[100vh] w-full max-w-[794px]'>
      <Invoice
        check={values?.check}
        availableBankAccount={displayedBankAccount}
        isInvoiceBankDetailsHovered={!!hoveredAccount}
        isInvoiceBankDetailsFocused={!!focusedAccount}
        canPayWithBankTransfer={values?.canPayWithBankTransfer}
        currency={displayedBankAccount?.currency || values?.baseCurrency?.value}
        isInvoiceProviderFocused={
          isInvoiceProviderFocused || isInvoiceProviderDetailsHovered
        }
        from={{
          addressLine1: values?.addressLine1 ?? '',
          addressLine2: values?.addressLine2 ?? '',
          locality: values?.locality ?? '',
          zip: values?.zip ?? '',
          country: values?.country?.label ?? '',
          region: values?.region ?? '',
          email: values?.sendInvoicesFrom ?? '',
          name: values?.legalName ?? '',
          vatNumber: values?.vatNumber ?? '',
        }}
        {...invoicePreviewStaticData}
      />
    </div>
  );
};
