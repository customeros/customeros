'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import React, { useRef, useMemo, useState } from 'react';

import { produce } from 'immer';
import { useDeepCompareEffect } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { Box } from '@ui/layout/Box';
import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Grid, GridItem } from '@ui/layout/Grid';
import { Heading } from '@ui/typography/Heading';
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
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';

import { BillingDetailsDto } from './BillingDetails.dto';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  organizationName: string;
}
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const ContractBillingDetailsModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  notes,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
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

  const [isBillingDetailsFocused, setIsBillingDetailsFocused] =
    useState<boolean>(false);

  const [isBillingDetailsHovered, setIsBillingDetailsHovered] =
    useState<boolean>(false);
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
            (contract) => contract.id === contractId,
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
  const defaultValues = new BillingDetailsDto({
    ...(data?.contract ?? {}),
    organizationLegalName:
      data?.contract?.organizationLegalName || organizationName,
  });

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'invoiceEmail') {
          return {
            ...next,
            values: {
              ...next.values,
              invoiceEmail: action.payload.value.split(' ').join('').trim(),
            },
          };
        }
      }

      return next;
    },
  });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

  const handleApplyChanges = () => {
    const payload = BillingDetailsDto.toPayload(state.values);

    updateContract.mutate({
      input: {
        contractId,
        ...payload,
      },
    });
  };
  const invoicePreviewStaticData = useMemo(
    () => ({
      status: 'Preview',
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
            email: '',
            name: tenantBillingProfile?.tenantBillingProfiles?.[0]?.legalName,
          }
        : {
            addressLine1: '29 Maple Lane',
            addressLine2: 'Springfield, Haven County',
            locality: 'San Francisco',
            zip: '89302',
            country: 'United States',
            email: 'invoices@acme.com',
            name: 'Acme Corp.',
          },
    }),
    [tenantBillingProfile?.tenantBillingProfiles?.[0]],
  );
  const isEmailValid = useMemo(() => {
    return (
      !!state.values.invoiceEmail?.length &&
      !emailRegex.test(state.values.invoiceEmail)
    );
  }, [state?.values?.invoiceEmail]);
  const availableCurrencies = useMemo(
    () => (bankAccountsData?.bankAccounts ?? []).map((e) => e.currency),
    [],
  );

  const canAllowPayWithBankTransfer = useMemo(() => {
    return availableCurrencies.includes(state.values.currency?.value);
  }, [availableCurrencies, state.values.currency]);

  useDeepCompareEffect(() => {
    if (!canAllowPayWithBankTransfer) {
      const newDefaultValues = new BillingDetailsDto({
        ...(data?.contract ?? {}),
        organizationLegalName:
          data?.contract?.organizationLegalName || organizationName,
        billingDetails: {
          ...(data?.contract?.billingDetails ?? {}),
          canPayWithBankTransfer: false,
        },
      });
      setDefaultValues(newDefaultValues);
    }
  }, [canAllowPayWithBankTransfer]);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='4xl'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <Grid h='100%' templateColumns='356px 1fr'>
          <GridItem
            rowSpan={1}
            colSpan={1}
            h='100%'
            display='flex'
            flexDir='column'
            justifyContent='space-between'
            bg='gray.25'
            borderRight='1px solid'
            borderColor='gray.200'
            borderTopLeftRadius='2xl'
            borderBottomLeftRadius='2xl'
            backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
            backgroundRepeat='no-repeat'
            sx={{
              backgroundPositionX: '1px',
              backgroundPositionY: '-7px',
            }}
          >
            <ModalHeader>
              <FeaturedIcon size='lg' colorScheme='primary'>
                <File02 color='primary.600' />
              </FeaturedIcon>
              <Heading fontSize='lg' mt='4'>
                {data?.contract?.organizationLegalName ||
                  organizationName ||
                  "Unnamed's "}{' '}
                contract details
              </Heading>
            </ModalHeader>
            <ContractBillingDetailsForm
              formId={formId}
              tenantBillingProfile={
                tenantBillingProfile
                  ?.tenantBillingProfiles?.[0] as TenantBillingProfile
              }
              organizationName={organizationName}
              canAllowPayWithBankTransfer={canAllowPayWithBankTransfer}
              hasNoBankAccounts={!bankAccountsData?.bankAccounts?.length}
              currency={state?.values?.currency?.value}
              isEmailValid={isEmailValid}
              onSetIsBillingDetailsHovered={setIsBillingDetailsHovered}
              onSetIsBillingDetailsFocused={setIsBillingDetailsFocused}
            />
            <ModalFooter p='6'>
              <Button variant='outline' w='full' onClick={onClose}>
                Cancel
              </Button>
              <Button
                ml='3'
                w='full'
                variant='outline'
                colorScheme='primary'
                loadingText='Applying changes...'
                onClick={handleApplyChanges}
              >
                Done
              </Button>
            </ModalFooter>
          </GridItem>
          <GridItem>
            <Box width='100%' h='full'>
              <Invoice
                isBilledToFocused={
                  isBillingDetailsFocused || isBillingDetailsHovered
                }
                note={notes}
                currency={state?.values?.currency?.value}
                billedTo={{
                  addressLine1: state.values.addressLine1 ?? '',
                  addressLine2: state.values.addressLine2 ?? '',
                  locality: state.values.locality ?? '',
                  zip: state.values.zip ?? '',
                  country: state?.values?.country?.label ?? '',
                  email: state.values.invoiceEmail ?? '',
                  name: state.values?.organizationLegalName ?? '',
                }}
                {...invoicePreviewStaticData}
                canPayWithBankTransfer={
                  tenantBillingProfile?.tenantBillingProfiles?.[0]
                    ?.canPayWithBankTransfer &&
                  state.values.canPayWithBankTransfer
                }
                availableBankAccount={
                  bankAccountsData?.bankAccounts?.find(
                    (e) => e.currency === state?.values?.currency?.value,
                  ) as BankAccount
                }
              />
            </Box>
          </GridItem>
        </Grid>
      </ModalContent>
    </Modal>
  );
};
