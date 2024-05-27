import { useForm } from 'react-inverted-form';
import { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { observer } from 'mobx-react-lite';
import { useQueryClient } from '@tanstack/react-query';

import { DateTimeUtils } from '@spaces/utils/date';
import { toastError } from '@ui/presentation/Toast';
import { FormInput } from '@ui/form/Input/FormInput';
import { Contract, ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardFooter, CardHeader } from '@ui/presentation/Card/Card';
import { useUpdateContractMutation } from '@organization/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import { useUpdatePanelModalStateContext } from '@organization/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { UpcomingInvoices } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/UpcomingInvoices/UpcomingInvoices';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { useEditContractModalStores } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/EditContractModalStores';

import { Services } from './Services/Services';
import { ContractSubtitle } from './ContractSubtitle';
import { ContractCardActions } from './ContractCardActions';
import { RenewalARRCard } from './RenewalARR/RenewalARRCard';
import { EditContractModal } from './ContractBillingDetailsModal/EditContractModal';

interface ContractCardProps {
  data: Contract;
  organizationId: string;
  organizationName: string;
}

export const ContractCard = observer(
  ({ data, organizationName, organizationId }: ContractCardProps) => {
    const queryKey = useGetContractsQuery.getKey({ id: organizationId });
    const { serviceFormStore } = useEditContractModalStores();

    const queryClient = useQueryClient();
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);
    const [isExpanded, setIsExpanded] = useState(!data?.contractSigned);
    const formId = `contract-form-${data.metadata.id}`;
    const { setIsPanelModalOpen } = useUpdatePanelModalStateContext();
    const {
      isEditModalOpen,
      onEditModalOpen,
      onChangeModalMode,
      onEditModalClose,
    } = useContractModalStateContext();

    const client = getGraphQLClient();

    // this is needed to block scroll on safari when modal is open, scrollbar overflow issue
    useEffect(() => {
      if (isEditModalOpen) {
        setIsPanelModalOpen(true);
      }
      if (!isEditModalOpen) {
        serviceFormStore.clearUsedColors();

        setIsPanelModalOpen(false);
      }
    }, [isEditModalOpen]);

    useEffect(() => {
      serviceFormStore.contractIdValue = data.metadata.id;
      if (data.contractLineItems?.length && isEditModalOpen) {
        serviceFormStore.initializeServices(data.contractLineItems);
      }
    }, [isEditModalOpen, data.contractLineItems]);

    const updateContract = useUpdateContractMutation(client, {
      onMutate: ({ input: { patch, contractId, ...input } }) => {
        queryClient.cancelQueries({ queryKey });
        queryClient.setQueryData<GetContractsQuery>(
          queryKey,
          (currentCache) => {
            return produce(currentCache, (draft) => {
              const previousContracts = draft?.['organization']?.['contracts'];
              const updatedContractIndex = previousContracts?.findIndex(
                (contract) => contract.metadata.id === data?.metadata?.id,
              );
              if (draft?.['organization']?.['contracts']) {
                draft['organization']['contracts']?.map(
                  (contractData, index) => {
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
                  },
                );
              }
            });
          },
        );
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
      stateReducer: (_state, action, next) => {
        if (
          action.type === 'FIELD_CHANGE' &&
          action.payload.name === 'contractName'
        ) {
          updateContractDebounced(
            action.payload.value.length === 0
              ? 'Unnamed contract'
              : action.payload.value,
          );
        }

        return next;
      },
    });

    const handleOpenBillingDetails = () => {
      onChangeModalMode(EditModalMode.BillingDetails);
      onEditModalOpen();
    };
    const handleOpenContractDetails = () => {
      onChangeModalMode(EditModalMode.ContractDetails);
      onEditModalOpen();
    };

    return (
      <Card className='px-4 py-3 w-full text-lg bg-gray-50 transition-all-0.2s-ease-out border border-gray-200 text-gray-700 '>
        <CardHeader
          className='p-0 w-full flex flex-col'
          role='button'
          onClick={() => (!isExpanded ? setIsExpanded(true) : null)}
        >
          <article className='flex justify-between flex-1 w-full'>
            <FormInput
              className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
              name='contractName'
              formId={formId}
            />

            <ContractCardActions
              onOpenEditModal={handleOpenContractDetails}
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
            />
          </article>

          <div
            role='button'
            tabIndex={1}
            onClick={handleOpenContractDetails}
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
            onModalOpen={onEditModalOpen}
          />
          {!!data?.upcomingInvoices?.length && (
            <>
              <Divider className='my-3' />
              <UpcomingInvoices
                data={data}
                onOpenBillingDetailsModal={handleOpenBillingDetails}
                onOpenServiceLineItemsModal={handleOpenContractDetails}
              />
            </>
          )}

          <EditContractModal
            isOpen={isEditModalOpen}
            status={data?.contractStatus}
            contractId={data.metadata.id}
            onClose={onEditModalClose}
            serviceStarted={data.serviceStarted}
            organizationName={organizationName}
            notes={data?.billingDetails?.invoiceNote}
            renewsAt={data?.opportunities?.[0]?.renewedAt}
          />
        </CardFooter>
      </Card>
    );
  },
);
