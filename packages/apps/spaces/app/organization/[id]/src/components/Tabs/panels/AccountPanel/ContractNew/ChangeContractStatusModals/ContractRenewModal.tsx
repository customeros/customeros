'use client';
import React, { useRef } from 'react';
import { useParams } from 'next/navigation';

import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { FeaturedIcon } from '@ui/media/Icon';
import { ContractStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { useRenewContractMutation } from '@organization/src/graphql/renewContract.generated';
import { useContractModalStatusContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractEndModalProps {
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  organizationName: string;

  status?: ContractStatus | null;
}

export const ContractRenewsModal = ({
  onClose,
  contractId,
  organizationName,
  status,
  renewsAt,
}: ContractEndModalProps) => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;

  const queryKey = useGetContractsQuery.getKey({ id });
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { nextInvoice } = useContractModalStatusContext();

  const { mutate } = useRenewContractMutation(client, {
    onSuccess: () => {
      onClose();
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({
          queryKey: ['GetTimeline.infinite'],
        });
      }, 1000);
    },
  });

  const getText = () => {
    if (
      status === ContractStatus.OutOfContract &&
      renewsAt &&
      DateTimeUtils.isSameDay(renewsAt, new Date().toISOString())
    ) {
      return 'Renewing this contract will extend it with another year until 2 Aug 2025';
    }
  };

  return (
    <>
      <div>
        {!nextInvoice && (
          <FeaturedIcon size='lg' colorScheme='primary'>
            <RefreshCw05 className='text-primary-600' />
          </FeaturedIcon>
        )}

        <h1
          className={cn('text-lg font-semibold  mb-1', {
            'mt-4': !nextInvoice,
          })}
        >
          {status === ContractStatus.OutOfContract
            ? 'Renew this contract?'
            : 'When should this contract renew?'}
        </h1>
      </div>

      <p className='flex flex-col gap-3'>{getText()}</p>
      <div className='flex'>
        <Button
          variant='outline'
          colorScheme='gray'
          className='w-full'
          onClick={onClose}
        >
          Cancel
        </Button>
        <Button
          className='ml-3 w-full'
          variant='outline'
          colorScheme='primary'
          onClick={() => mutate({ contractId })}
        >
          Renew now
        </Button>
      </div>
    </>
  );
};
