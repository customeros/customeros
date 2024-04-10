'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import React, { useRef, useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useDeepCompareEffect } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { Button } from '@ui/form/Button/Button';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useGetContractQuery } from '@organization/src/graphql/getContract.generated';
import { useUpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  DataSource,
  BankAccount,
  InvoiceLine,
  TenantBillingProfile,
} from '@graphql/types';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';

import { ContractDetailsDto } from './ContractDetails.dto';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  organizationName: string;
  billedTo: {
    zip: string;
    name: string;
    email: string;
    region: string;
    country: string;
    locality: string;
    addressLine1: string;
    addressLine2: string;
  };
}

export const EditContractModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  notes,
  billedTo,
}: SubscriptionServiceModalProps) => {
  const formId = `billing-details-form-${contractId}`;
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();

  const { data } = useGetContractQuery(
    client,
    {
      id: contractId,
    },
    {
      enabled: isOpen && !!contractId,
      refetchOnMount: true,
    },
  );
  const { data: bankAccountsData } = useBankAccountsQuery(client);

  const queryKey = useGetContractsQuery.getKey({ id: organizationId });

  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { data: tenantBillingProfile } = useTenantBillingProfilesQuery(client);

  const updateContract = useUpdateContractMutation(client, {
    onMutate: ({
      input: {
        patch,
        contractId,
        canPayWithBankTransfer,
        canPayWithDirectDebit,
        canPayWithCard,
        ...input
      },
    }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.metadata.id === contractId,
          );

          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return {
                ...contractData,
                ...input,
              };
            });
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (error, _, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );

      toastError(
        'Failed to update billing details',
        `update-contract-error-${error}`,
      );
    },
    onSuccess: () => {
      toastSuccess(
        'Billing details updated',
        `update-contract-success-${contractId}`,
      );
      onClose();
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
      }, 500);
    },
  });
  const defaultValues = new ContractDetailsDto(data?.contract);

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

  const handleApplyChanges = () => {
    const payload = ContractDetailsDto.toPayload(state.values);
    updateContract.mutate({
      input: {
        contractId,
        ...payload,
      },
    });
  };
  const invoicePreviewStaticData = useMemo(
    () => ({
      status: null,
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
        } as InvoiceLine,
      ],
      tax: 0,
      total: 100,
      dueDate: new Date().toISOString(),
      subtotal: 100,
      issueDate: new Date().toISOString(),
      from: tenantBillingProfile?.tenantBillingProfiles?.[0]
        ? {
            addressLine1:
              tenantBillingProfile?.tenantBillingProfiles?.[0]?.addressLine1 ??
              '',
            addressLine2:
              tenantBillingProfile?.tenantBillingProfiles?.[0].addressLine2,
            locality:
              tenantBillingProfile?.tenantBillingProfiles?.[0]?.locality ?? '',
            zip: tenantBillingProfile?.tenantBillingProfiles?.[0]?.zip ?? '',
            country: tenantBillingProfile?.tenantBillingProfiles?.[0].country
              ? countryOptions.find(
                  (country) =>
                    country.value ===
                    tenantBillingProfile?.tenantBillingProfiles?.[0]?.country,
                )?.label
              : '',
            email:
              tenantBillingProfile?.tenantBillingProfiles?.[0]
                ?.sendInvoicesFrom,
            name: tenantBillingProfile?.tenantBillingProfiles?.[0]?.legalName,
            region:
              tenantBillingProfile?.tenantBillingProfiles?.[0]?.region ?? '',
          }
        : {
            addressLine1: '29 Maple Lane',
            addressLine2: 'Springfield, Haven County',
            locality: 'San Francisco',
            zip: '89302',
            country: 'United States of America',
            email: 'invoices@acme.com',
            name: 'Acme Corp.',
            region: 'California',
          },
    }),
    [tenantBillingProfile?.tenantBillingProfiles?.[0]],
  );

  const availableCurrencies = useMemo(
    () => (bankAccountsData?.bankAccounts ?? []).map((e) => e.currency),
    [],
  );

  const canAllowPayWithBankTransfer = useMemo(() => {
    return availableCurrencies.includes(state.values.currency?.value);
  }, [availableCurrencies, state.values.currency]);

  useEffect(() => {
    if (!canAllowPayWithBankTransfer) {
      const newDefaultValues = new ContractDetailsDto({
        ...(data?.contract ?? {}),
        billingDetails: {
          ...(data?.contract?.billingDetails ?? {}),
          canPayWithBankTransfer: false,
        },
      });
      setDefaultValues(newDefaultValues);
    }
  }, [canAllowPayWithBankTransfer]);

  const handleCloseModal = () => {
    setDefaultValues(defaultValues);
    onClose();
  };

  return (
    <Modal open={isOpen} onOpenChange={onClose}>
      <ModalOverlay />
      <ModalContent
        className='border-r-2 flex gap-6 bg-transparent shadow-none border-none'
        style={{ minWidth: '971px', minHeight: '80%', boxShadow: 'none' }}
      >
        <div
          className='flex flex-col gap-4 px-6 pb-6 pt-4 bg-white h-auto rounded-lg justify-between'
          style={{ width: '424px' }}
        >
          <ModalHeader className='p-0 text-lg font-semibold'>
            {data?.contract?.organizationLegalName ||
              organizationName ||
              "Unnamed's "}{' '}
            contract details
          </ModalHeader>
          <ContractBillingDetailsForm
            formId={formId}
            contractId={contractId}
            tenantBillingProfile={
              tenantBillingProfile
                ?.tenantBillingProfiles?.[0] as TenantBillingProfile
            }
            currency={state?.values?.currency?.value}
            bankAccounts={bankAccountsData?.bankAccounts as BankAccount[]}
            payAutomatically={state?.values?.payAutomatically}
          />
          <ModalFooter className='p-0 flex'>
            <Button
              variant='outline'
              colorScheme='gray'
              onClick={handleCloseModal}
              className='w-full'
              size='md'
            >
              Cancel
            </Button>
            <Button
              className='ml-3 w-full'
              size='md'
              variant='outline'
              colorScheme='primary'
              onClick={handleApplyChanges}
            >
              Confirm
            </Button>
          </ModalFooter>
        </div>
        <div style={{ width: '570px' }} className='bg-white rounded'>
          <div className='w-full h-full'>
            <Invoice
              isBilledToFocused={false}
              note={notes}
              currency={state?.values?.currency?.value}
              billedTo={billedTo}
              {...invoicePreviewStaticData}
              canPayWithBankTransfer={
                tenantBillingProfile?.tenantBillingProfiles?.[0]
                  ?.canPayWithBankTransfer &&
                state.values.canPayWithBankTransfer
              }
              check={tenantBillingProfile?.tenantBillingProfiles?.[0]?.check}
              availableBankAccount={
                bankAccountsData?.bankAccounts?.find(
                  (e) => e.currency === state?.values?.currency?.value,
                ) as BankAccount
              }
            />
          </div>
        </div>
      </ModalContent>
    </Modal>
  );
};
