import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
import { BankAccount, InvoiceLine, InvoiceStatus } from '@graphql/types';
import { ISimulatedInvoiceLineItems } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/InvoicePreviewList.store.ts';

import { ServicesTable } from './ServicesTable';
import logoCustomerOs from './assets/customer-os.png';
import {
  InvoiceHeader,
  InvoiceSummary,
  BankingDetails,
  InvoicePartySection,
} from './components';

// todo refactor, use generated type
type Address = {
  zip: string;
  email: string;
  name?: string;
  region: string;
  country?: string;
  locality: string;
  vatNumber?: string;
  addressLine1: string;
  addressLine2?: string;
};

type InvoiceProps = {
  tax: number;
  from: Address;
  total: number;
  dueDate: string;
  subtotal: number;
  currency?: string;
  issueDate: string;
  billedTo: Address;
  amountDue?: number;
  note?: string | null;
  invoiceNumber: string;
  check?: boolean | null;
  invoicePeriodEnd?: string;
  shouldBlurDummy?: boolean;
  isBilledToFocused?: boolean;
  invoicePeriodStart?: string;
  status?: InvoiceStatus | null;
  isInvoiceProviderFocused?: boolean;
  isInvoiceBankDetailsHovered?: boolean;
  isInvoiceBankDetailsFocused?: boolean;
  onOpenAddressDetailsModal?: () => void;
  canPayWithBankTransfer?: boolean | null;
  availableBankAccount?: Partial<BankAccount> | null;
  lines: InvoiceLine[] | ISimulatedInvoiceLineItems[];
};

export function Invoice({
  invoiceNumber,
  issueDate,
  dueDate,
  billedTo,
  from,
  lines,
  subtotal,
  tax,
  total,
  note,
  amountDue,
  status,
  isBilledToFocused,
  isInvoiceProviderFocused,
  isInvoiceBankDetailsHovered,
  isInvoiceBankDetailsFocused,
  currency = 'USD',
  canPayWithBankTransfer,
  availableBankAccount,
  check,
  invoicePeriodStart,
  invoicePeriodEnd,
  onOpenAddressDetailsModal,
  shouldBlurDummy,
}: InvoiceProps) {
  const isInvoiceMetaSectionBlurred =
    isBilledToFocused || isInvoiceProviderFocused;

  const invoiceMetaSectionFilterProperty = isInvoiceMetaSectionBlurred
    ? 'blur-[2px]'
    : 'filter-none';

  const isInvoiceBankDetailsSectionFocused =
    (isInvoiceBankDetailsHovered && !isInvoiceMetaSectionBlurred) ||
    isInvoiceBankDetailsFocused;

  const isInvoiceTopSectionFilterProperty = isInvoiceBankDetailsSectionFocused
    ? 'blur-[2px]'
    : 'filter-none';

  const blurDummyClass = shouldBlurDummy ? 'blur-[2px]' : 'filter-none';

  return (
    <div className='px-4 flex flex-col w-full overflow-y-auto h-full justify-between pb-4 '>
      <div className={cn('flex flex-col', isInvoiceTopSectionFilterProperty)}>
        <div className='flex relative flex-col mt-2'>
          <InvoiceHeader invoiceNumber={invoiceNumber} status={status} />

          <div className='flex mt-2 justify-evenly transition duration-250 ease-in-out filter'>
            <div
              className={cn(
                'flex flex-1 flex-col w-170 py-2 px-2 border-r border-t border-b border-gray-300 transition duration-250 ease-in-out filter',
                invoiceMetaSectionFilterProperty,
              )}
            >
              <span className='font-semibold mb-1 text-sm'>Issued</span>
              <span
                className={cn('text-sm mb-4 text-gray-500', blurDummyClass)}
              >
                {DateTimeUtils.format(
                  issueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </span>
              <span className='font-semibold mb-1 text-sm'>Due</span>
              <span className={cn('text-sm text-gray-500', blurDummyClass)}>
                {DateTimeUtils.format(
                  dueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </span>
            </div>
            <InvoicePartySection
              title='Billed to'
              isBlurred={isInvoiceProviderFocused}
              isFocused={isBilledToFocused}
              zip={billedTo?.zip}
              name={billedTo?.name}
              email={billedTo?.email}
              country={billedTo?.country}
              locality={billedTo.locality}
              addressLine1={billedTo?.addressLine1}
              addressLine2={billedTo?.addressLine2}
              vatNumber={billedTo?.vatNumber}
              region={billedTo?.region}
              onClick={onOpenAddressDetailsModal}
            />
            <InvoicePartySection
              title='From'
              isBlurred={isBilledToFocused}
              isFocused={isInvoiceProviderFocused}
              zip={from?.zip}
              name={from?.name}
              email={from?.email}
              country={from?.country}
              locality={from?.locality}
              addressLine1={from?.addressLine1}
              addressLine2={from?.addressLine2}
              vatNumber={from?.vatNumber}
              region={from?.region}
            />
          </div>
        </div>

        <div
          className={cn(
            'flex flex-col mt-4 transition duration-250 ease-in-out filter',
            isInvoiceTopSectionFilterProperty,
          )}
        >
          <ServicesTable
            services={lines ?? []}
            currency={currency}
            invoicePeriodEnd={invoicePeriodEnd}
            invoicePeriodStart={invoicePeriodStart}
          />
          <InvoiceSummary
            tax={tax}
            total={total}
            subtotal={subtotal}
            currency={currency}
            amountDue={amountDue}
            note={note}
          />
        </div>
      </div>

      <div
        className={isInvoiceMetaSectionBlurred ? 'filter-[2px]' : 'filter-none'}
      >
        <div
          className={cn('w-full, border-y-2', {
            'border-gray-900': isInvoiceBankDetailsSectionFocused,
            'border-transparent': !isInvoiceBankDetailsSectionFocused,
          })}
        >
          {canPayWithBankTransfer && availableBankAccount && (
            <BankingDetails
              availableBankAccount={availableBankAccount}
              currency={currency}
              invoiceNumber={invoiceNumber}
            />
          )}
        </div>

        {check && (
          <span className='text-xs text-gray-500 my-2'>
            Want to pay by check? Contact{' '}
            <a
              className='underline text-gray-500'
              href={`mailto:${from.email}`}
            >
              {from.email}
            </a>
          </span>
        )}
        <div className='flex items-center py-2 mt-2 border-t border-gray-300'>
          <div className='mr-2'>
            <img src={logoCustomerOs} alt='CustomerOS' width={14} height={14} />
          </div>
          <span className='text-xs text-gray-500'>
            Powered by
            <a
              className='text-gray-500 mx-1 underline cursor-pointer'
              href='https://customeros.ai/'
              target='blank'
              rel='noopener noreferrer'
            >
              CustomerOS
            </a>
            - Revenue Intelligence for B2B hyperscalers
          </span>
        </div>
      </div>
    </div>
  );
}
