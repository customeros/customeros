import { useForm } from 'react-inverted-form';
import React, { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';

import { useDisclosure } from '@ui/utils';
import { DateTimeUtils } from '@spaces/utils/date';
import { toastError } from '@ui/presentation/Toast';
import { FormInput } from '@ui/form/Input/FormInput2';
import { Contract, ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardFooter, CardHeader } from '@ui/presentation/Card/Card';
import { useUpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { useUpdatePanelModalStateContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { UpcomingInvoices } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/UpcomingInvoices/UpcomingInvoices';

import { Services } from './Services/Services';
import { ContractSubtitle } from './ContractSubtitle';
import { ContractCardActions } from './ContractCardActions';
import { RenewalARRCard } from './RenewalARR/RenewalARRCard';
import { ServiceLineItemsModal } from './ServiceLineItemsModal';
import { EditContractModal } from './ContractBillingDetailsModal/EditContractModal';

interface ContractCardProps {
  data: Contract;
  organizationId: string;
  organizationName: string;
}

export const ContractCard = ({
  data,
  organizationName,
  organizationId,
}: ContractCardProps) => {
  const queryKey = useGetContractsQuery.getKey({ id: organizationId });
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isExpanded, setIsExpanded] = useState(!data?.contractSigned);
  const formId = `contract-form-${data.metadata.id}`;
  const { setIsPanelModalOpen } = useUpdatePanelModalStateContext();
  const [isEditModalOpen, setEditModalOpen] = useState(false);
  const {
    onOpen: onServiceLineItemsOpen,
    onClose: onServiceLineItemClose,
    isOpen: isServceItemsModalOpen,
  } = useDisclosure({
    id: 'service-line-items-modal',
  });

  const client = getGraphQLClient();

  // this is needed to block scroll on safari when modal is open, scrollbar overflow issue
  useEffect(() => {
    if (isEditModalOpen || isServceItemsModalOpen) {
      setIsPanelModalOpen(true);
    }
    if (!isEditModalOpen && !isServceItemsModalOpen) {
      setIsPanelModalOpen(false);
    }
  }, [isEditModalOpen, isServceItemsModalOpen]);

  const updateContract = useUpdateContractMutation(client, {
    onMutate: ({ input: { patch, contractId, ...input } }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.metadata.id === data?.metadata?.id,
          );
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }
              const result = Object.entries(input).find(
                ([_, value]) => value === '0001-01-01T00:00:00.000000Z',
              );

              return {
                ...contractData,
                ...input,
                ...(result ? { [result[0]]: null } : {}),
              };
            });
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (error, { input }, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );

      const invalidDate =
        DateTimeUtils.isBefore(input.contractEnded, input.serviceStarted) ||
        DateTimeUtils.isBefore(input.contractEnded, input.contractSigned);

      toastError(
        `${
          invalidDate
            ? 'The contract must end after the service start or signing date'
            : 'Failed to update contract'
        }`,
        `update-contract-error-${error}`,
      );
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });

        queryClient.invalidateQueries({ queryKey: ['GetTimeline.infinite'] });
      }, 1000);
    },
  });
  const updateContractDebounced = useDebounce((name: string) => {
    updateContract.mutate({
      input: {
        contractId: data.metadata.id,
        contractName: name,
        patch: true,
      },
    });
  }, 500);

  useForm<{ contractName: string }>({
    formId,
    defaultValues: {
      contractName: data.contractName,
    },
    debug: true,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'name') {
          updateContractDebounced(action.payload.value);
        }
      }

      return next;
    },
  });

  return (
    <Card className='px-4 py-3 w-full text-lg bg-gray-50 transition-all-0.2s-ease-out border border-gray-200 '>
      <CardHeader
        className='p-0 w-full flex flex-col'
        role='button'
        onClick={() => (!isExpanded ? setIsExpanded(true) : null)}
      >
        <article className='flex justify-between flex-1 w-full'>
          <FormInput
            className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0'
            name='contractName'
            formId={formId}
          />

          <ContractCardActions
            onOpenEditModal={() => setEditModalOpen(true)}
            status={data.contractStatus}
            contractId={data.metadata.id}
            renewsAt={data?.opportunities?.[0]?.renewedAt}
            onUpdateContract={updateContract}
            serviceStarted={data.serviceStarted}
            organizationName={
              data?.billingDetails?.organizationLegalName ||
              organizationName ||
              'Unnamed'
            }
            nextInvoiceDate={data?.billingDetails?.nextInvoicing}
            contractStarted={data.serviceStarted}
          />
        </article>

        <div
          role='button'
          tabIndex={1}
          onClick={() => setEditModalOpen(true)}
          className='w-full'
        >
          <ContractSubtitle data={data} />
        </div>
      </CardHeader>

      <CardFooter className='p-0 mt-0 w-full flex flex-col'>
        {data?.opportunities && !!data.contractLineItems?.length && (
          <RenewalARRCard
            hasEnded={data.contractStatus === ContractStatus.Ended}
            startedAt={data.serviceStarted}
            currency={data.currency}
            opportunity={data.opportunities?.[0]}
          />
        )}
        <Services
          data={data?.contractLineItems}
          currency={data?.currency}
          onModalOpen={onServiceLineItemsOpen}
        />
        <Divider className='my-3' />

        <UpcomingInvoices data={data} />
        <EditContractModal
          isOpen={isEditModalOpen}
          contractId={data.metadata.id}
          onClose={() => setEditModalOpen(false)}
          organizationName={organizationName}
          notes={data?.billingDetails?.invoiceNote}
        />

        <ServiceLineItemsModal
          isOpen={isServceItemsModalOpen}
          contractId={data.metadata.id}
          onClose={onServiceLineItemClose}
          contractName={data.contractName}
          currency={data.currency}
          contractLineItems={data?.contractLineItems ?? []}
          organizationName={organizationName}
          notes={data?.billingDetails?.invoiceNote}
        />
      </CardFooter>
    </Card>
  );
};
