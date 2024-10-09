import { toZonedTime } from 'date-fns-tz';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import {
  BankAccount,
  InvoiceLine,
  InvoiceStatus,
  InvoiceLineSimulate,
} from '@graphql/types';

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
  billingPeriodsInMonths?: number | null;
  canPayWithBankTransfer?: boolean | null;
  lines: InvoiceLine[] | InvoiceLineSimulate[];
  availableBankAccount?: Partial<BankAccount> | null;
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
  billingPeriodsInMonths,
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
  const utcDueDate = toZonedTime(dueDate, 'UTC').toUTCString();
  const utcIssueDate = toZonedTime(issueDate, 'UTC').toUTCString();

  return (
    <div className='px-4 flex flex-col w-full overflow-y-auto h-full justify-between pb-4 '>
      <div className={cn('flex flex-col', isInvoiceTopSectionFilterProperty)}>
        <div className='flex relative flex-col mt-2'>
          <InvoiceHeader status={status} invoiceNumber={invoiceNumber} />

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
                  utcIssueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </span>
              <span className='font-semibold mb-1 text-sm'>Due</span>
              <span className={cn('text-sm text-gray-500', blurDummyClass)}>
                {DateTimeUtils.format(
                  utcDueDate,
                  DateTimeUtils.dateWithAbreviatedMonth,
                )}
              </span>
            </div>
            <InvoicePartySection
              title='Billed to'
              zip={billedTo?.zip}
              name={billedTo?.name}
              email={billedTo?.email}
              region={billedTo?.region}
              country={billedTo?.country}
              locality={billedTo.locality}
              isFocused={isBilledToFocused}
              vatNumber={billedTo?.vatNumber}
              onClick={onOpenAddressDetailsModal}
              isBlurred={isInvoiceProviderFocused}
              addressLine1={billedTo?.addressLine1}
              addressLine2={billedTo?.addressLine2}
            />
            <InvoicePartySection
              title='From'
              zip={from?.zip}
              name={from?.name}
              email={from?.email}
              region={from?.region}
              country={from?.country}
              locality={from?.locality}
              vatNumber={from?.vatNumber}
              isBlurred={isBilledToFocused}
              addressLine1={from?.addressLine1}
              addressLine2={from?.addressLine2}
              isFocused={isInvoiceProviderFocused}
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
            currency={currency}
            services={lines ?? []}
            invoicePeriodEnd={invoicePeriodEnd}
            invoicePeriodStart={invoicePeriodStart}
            billingPeriodsInMonths={billingPeriodsInMonths}
          />
          <InvoiceSummary
            tax={tax}
            note={note}
            total={total}
            subtotal={subtotal}
            currency={currency}
            amountDue={amountDue}
          />
        </div>
      </div>

      <div
        className={isInvoiceMetaSectionBlurred ? 'filter-[2px]' : 'filter-none'}
      >
        <div
          className={cn('w-full, border-y-2 mt-10', {
            'border-gray-900': isInvoiceBankDetailsSectionFocused,
            'border-transparent': !isInvoiceBankDetailsSectionFocused,
          })}
        >
          {canPayWithBankTransfer && availableBankAccount && (
            <BankingDetails
              currency={currency}
              invoiceNumber={invoiceNumber}
              availableBankAccount={availableBankAccount}
            />
          )}
        </div>

        {check && (
          <div className='text-xs text-gray-500 my-2 pt-2 border-t border-gray-300'>
            Want to pay by check? Contact{' '}
            <a
              href={`mailto:${from.email}`}
              className='underline text-gray-500'
            >
              {from.email}
            </a>
          </div>
        )}
        <div className='flex items-center py-2 mt-2 border-t border-gray-300'>
          <div className='mr-2'>
            <img width={14} height={14} alt='CustomerOS' src={logoCustomerOs} />
          </div>
          <span className='text-xs text-gray-500'>
            Powered by
            <a
              target='blank'
              rel='noopener noreferrer'
              href='https://customeros.ai/'
              className='text-gray-500 mx-1 underline cursor-pointer'
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
