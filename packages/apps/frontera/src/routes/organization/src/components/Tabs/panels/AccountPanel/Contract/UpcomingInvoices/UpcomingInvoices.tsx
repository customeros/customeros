import { useState } from 'react';

import { useDeepCompareEffect } from 'rooks';

import { Play } from '@ui/media/icons/Play';
import { Plus } from '@ui/media/icons/Plus';
import { Edit03 } from '@ui/media/icons/Edit03';
import { Button } from '@ui/form/Button/Button';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Invoice, Contract, ContractStatus } from '@graphql/types';
import { ArrowNarrowRight } from '@ui/media/icons/ArrowNarrowRight';
import { UpcomingInvoice } from '@organization/components/Tabs/panels/AccountPanel/Contract/UpcomingInvoices/UpcomingInvoice.tsx';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

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

  const { onStatusModalOpen } = useContractModalStatusContext();

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
          size='xxs'
          colorScheme='warning'
          onClick={onOpenBillingDetailsModal}
          leftIcon={<Edit03 className='size-3' />}
        >
          Complete billing details
        </Button>
      );
    }
    if (!data?.contractLineItems?.length) {
      return (
        <Button
          className='ml-2 font-normal rounded'
          size='xxs'
          colorScheme='primary'
          onClick={onOpenServiceLineItemsModal}
          leftIcon={<Plus className='size-3' />}
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
          onClick={() => onStatusModalOpen(ContractStatusModalMode.Renew)}
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
          size='xxs'
          colorScheme='primary'
          onClick={() => onStatusModalOpen(ContractStatusModalMode.Start)}
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
        {data?.upcomingInvoices.map((invoice: Invoice) => (
          <UpcomingInvoice
            key={invoice?.metadata?.id}
            id={invoice?.metadata?.id}
          />
        ))}
      </div>
    </article>
  );
};
