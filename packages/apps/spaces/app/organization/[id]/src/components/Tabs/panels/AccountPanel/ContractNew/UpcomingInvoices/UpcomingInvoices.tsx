import React, { useState } from 'react';

import { useDeepCompareEffect } from 'rooks';

import { Play } from '@ui/media/icons/Play';
import { Plus } from '@ui/media/icons/Plus';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { toastSuccess } from '@ui/presentation/Toast';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Contract, ContractStatus } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ArrowNarrowRight } from '@ui/media/icons/ArrowNarrowRight';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useRenewContractMutation } from '@organization/src/graphql/renewContract.generated';

interface ContractCardProps {
  data: Contract;
  onOpenBillingDetailsModal: () => void;
  onOpenServiceLineItemsModal: () => void;
}

export const UpcomingInvoices = ({
  data,
  onOpenBillingDetailsModal,
  onOpenServiceLineItemsModal,
}: ContractCardProps) => {
  const [isPaused, setIsPaused] = useState(false);
  const [isMissingFields, setFieldsMissing] = useState(false);
  const client = getGraphQLClient();

  const { mutate: renewContract } = useRenewContractMutation(client, {
    onSuccess: () => {
      toastSuccess('Contract renewed', `${data?.metadata?.id}-renewed`);
    },
  });
  const getIsPaused = (): boolean => {
    if (
      [
        ContractStatus.OutOfContract,
        ContractStatus.Draft,
        ContractStatus.Ended,
      ].includes(data.contractStatus)
    ) {
      return true;
    }
    const hasAllRequiredFields = [
      data?.billingDetails?.addressLine1,
      data?.billingDetails?.postalCode,
      data?.billingDetails?.locality,
      data?.billingDetails?.organizationLegalName,
    ].every((field) => !!field);

    if (!hasAllRequiredFields) {
      setFieldsMissing(true);

      return true;
    }

    return !data?.contractLineItems?.length;
  };
  useDeepCompareEffect(() => {
    const paused = getIsPaused();
    setIsPaused(paused);
  }, [data]);

  const getActionButton = () => {
    if (isMissingFields) {
      return (
        <Button
          className='ml-2 font-normal rounded'
          size='xs'
          colorScheme='warning'
          onClick={onOpenBillingDetailsModal}
          leftIcon={<Edit03 />}
        >
          Complete billing details
        </Button>
      );
    }
    if (!data?.contractLineItems?.length) {
      return (
        <Button
          className='ml-2 font-normal rounded'
          size='xs'
          colorScheme='primary'
          onClick={onOpenServiceLineItemsModal}
          leftIcon={<Plus />}
        >
          Add a service
        </Button>
      );
    }

    if (data.contractStatus === ContractStatus.OutOfContract) {
      return (
        <Button
          className='ml-2 font-normal rounded'
          size='xs'
          colorScheme='primary'
          onClick={() => renewContract({ contractId: data.metadata.id })}
          leftIcon={<RefreshCw05 />}
        >
          Renew contract
        </Button>
      );
    }
    if (data.contractStatus === ContractStatus.Draft) {
      return (
        <Button
          className='ml-2 font-normal rounded'
          size='xs'
          colorScheme='primary'
          onClick={() => {}}
          leftIcon={<Play />}
        >
          Make contract live
        </Button>
      );
    }
  };

  return (
    <article className='w-full'>
      <p className='text-sm font-semibold mb-1 flex'>
        <span className='whitespace-nowrap'>Next invoice</span>
        {isPaused && (
          <div className='flex w-full justify-between'>
            <div>
              <ArrowNarrowRight className='mx-1' />
              <span className='font-normal'>Paused</span>
            </div>

            {getActionButton()}
          </div>
        )}
      </p>
      <div>
        {data?.upcomingInvoices.map((invoice) => (
          <div key={invoice.metadata.id} className='flex  text-sm'>
            <div className='whitespace-nowrap mr-1'>Monthly recurring:</div>
            <div className='whitespace-nowrap text-gray-500 underline'>
              {formatCurrency(invoice.amountDue, 2, invoice.currency)} on{' '}
              {DateTimeUtils.format(
                invoice.due,
                DateTimeUtils.defaultFormatShortString,
              )}{' '}
              (
              {DateTimeUtils.format(
                invoice.invoicePeriodStart,
                DateTimeUtils.dateDayAndMonth,
              )}{' '}
              -{' '}
              {DateTimeUtils.format(
                invoice.invoicePeriodEnd,
                DateTimeUtils.dateDayAndMonth,
              )}
              )
            </div>
          </div>
        ))}
      </div>
    </article>
  );
};
