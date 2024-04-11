'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import React, { useRef, useMemo, useState, useEffect } from 'react';

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
import { BillingDetailsForm } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/BillingAddressDetails/BillingAddressDetailsForm';
import {
  BillingDetailsDto,
  BillingAddressDetailsFormDto,
} from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/BillingAddressDetails/BillingAddressDetailsForm.dto';

import { ContractDetailsDto } from './ContractDetails.dto';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  organizationName: string;
}

export const EditContractModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  notes,
}: SubscriptionServiceModalProps) => {
  const formId = `billing-details-form-${contractId}`;
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const [billingDetailsFormOpen, setBillingDetailsOpen] =
    useState<boolean>(false);
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

  const addressDetailsDefailtValues = new BillingDetailsDto(data?.contract);

  const { state: addressState, setDefaultValues: setDefaultAddressValues } =
    useForm<BillingAddressDetailsFormDto>({
      formId: 'billing-details-address-form',
      defaultValues: addressDetailsDefailtValues,
      stateReducer: (_, action, next) => {
        return next;
      },
    });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);
  useDeepCompareEffect(() => {
    setDefaultAddressValues(addressDetailsDefailtValues);
  }, [addressDetailsDefailtValues]);
  const handleCloseModal = () => {
    setDefaultValues(defaultValues);
    setDefaultAddressValues(addressDetailsDefailtValues);
    setBillingDetailsOpen(false);
    onClose();
  };

  const handleApplyChanges = () => {
    const payload = ContractDetailsDto.toPayload(state.values);
    updateContract.mutate(
      {
        input: {
          contractId,
          ...payload,
        },
      },
      {
        onSuccess: () => {
          toastSuccess(
            'Billing details updated',
            `update-contract-success-${contractId}`,
          );
          handleCloseModal();
        },
      },
    );
  };

  const handleSaveAddressChanges = () => {
    const payload = BillingDetailsDto.toPayload(addressState.values);
    updateContract.mutate(
      {
        input: {
          contractId,
          ...payload,
        },
      },
      {
        onSuccess: () => {
          setBillingDetailsOpen(false);
        },
      },
    );
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

  return (
    <Modal open={isOpen} onOpenChange={handleCloseModal}>
      <ModalOverlay />
      <ModalContent
        className='border-r-2 flex gap-6 bg-transparent shadow-none border-none'
        style={{ minWidth: '971px', minHeight: '80%', boxShadow: 'none' }}
      >
        <div
          className='flex flex-col gap-4 px-6 pb-6 pt-4 bg-white h-auto rounded-lg justify-between relative'
          style={{ width: '424px', maxWidth: '424px' }}
        >
          {billingDetailsFormOpen && (
            <>
              <div
                className='h-[95%] rounded-l-lg left-[-12px] bg-white absolute w-[12px] shadow-inner'
                style={{ boxShadow: 'inset -3px 1px 6px 0px #10182814' }}
              />

              <div className='flex flex-col relative'>
                <ModalHeader className='p-0 text-lg font-semibold'>
                  <div>
                    {data?.contract?.organizationLegalName ||
                      organizationName ||
                      "Unnamed's "}{' '}
                  </div>
                  <span className='text-sm font-normal'>
                    These details are required to issue invoices
                  </span>
                </ModalHeader>

                <BillingDetailsForm
                  values={addressState.values}
                  formId={'billing-details-address-form'}
                />
              </div>
              <ModalFooter className='p-0 flex'>
                <Button
                  variant='outline'
                  colorScheme='gray'
                  onClick={() => setBillingDetailsOpen(false)}
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
                  onClick={handleSaveAddressChanges}
                >
                  Confirm
                </Button>
              </ModalFooter>
            </>
          )}

          {!billingDetailsFormOpen && (
            <>
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
            </>
          )}
        </div>
        <div style={{ minWidth: '600px' }} className='bg-white rounded'>
          <div className='w-full h-full'>
            <Invoice
              onOpenAddressDetailsModal={() => setBillingDetailsOpen(true)}
              isBilledToFocused={false}
              note={notes}
              currency={state?.values?.currency?.value}
              billedTo={{
                addressLine1: addressState.values.addressLine1 ?? '',
                addressLine2: addressState.values.addressLine2 ?? '',
                locality: addressState.values.locality ?? '',
                zip: addressState.values.postalCode ?? '',
                country: addressState?.values?.country?.label ?? '',
                email: addressState.values.billingEmail ?? '',
                name: addressState.values?.organizationLegalName ?? '',
                region: addressState.values?.region ?? '',
              }}
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
