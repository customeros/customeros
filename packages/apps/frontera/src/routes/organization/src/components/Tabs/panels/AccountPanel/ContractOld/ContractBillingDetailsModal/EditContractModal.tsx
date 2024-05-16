import { useParams } from 'react-router-dom';
import { useForm } from 'react-inverted-form';
import { useRef, useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useDeepCompareEffect } from 'rooks';
import { motion, Variants } from 'framer-motion';
import { useQueryClient } from '@tanstack/react-query';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';

import { cn } from '@ui/utils/cn';
import { FormInput } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { useGetContractQuery } from '@organization/graphql/getContract.generated';
import { useUpdateContractMutation } from '@organization/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import {
  DataSource,
  BankAccount,
  InvoiceLine,
  ContractStatus,
  TenantBillingProfile,
} from '@graphql/types';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { BillingDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/ContractOld/BillingAddressDetails/BillingAddressDetailsForm';
import {
  BillingDetailsDto,
  BillingAddressDetailsFormDto,
} from '@organization/components/Tabs/panels/AccountPanel/ContractOld/BillingAddressDetails/BillingAddressDetailsForm.dto';

import { ContractDetailsDto } from './ContractDetails.dto';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  organizationName: string;
}

const mainVariants = {
  open: {
    w: '424px',
    minW: '424px',
    x: 0,
    position: 'relative',
    scale: 1,
    transition: { duration: 0.2, ease: 'easeOut' },
  },
  closed: {
    w: '100px',
    minW: '100px',
    position: 'absolute',
    x: '-32px',
    scale: 0.95,
    transition: { duration: 0.2, ease: 'easeOut' },
  },
};

const variants = {
  open: {
    opacity: 1,
    w: '424px',
    minW: '424px',
    x: 0,
    display: 'block',
    transition: { duration: 0.2, ease: 'easeOut' },
  },
  closed: {
    opacity: 0,
    w: '0px',
    minW: '0px',
    x: '100%',
    display: 'none',
    transition: { duration: 0.2, ease: 'easeOut' },
  },
};

export const EditContractModal = ({
  contractId,
  organizationName,
  notes,
  renewsAt,
}: SubscriptionServiceModalProps) => {
  const formId = `billing-details-form-${contractId}`;
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const [initialOpen, setInitialOpen] = useState(EditModalMode.ContractDetails);
  useState<boolean>(false);
  const {
    isEditModalOpen,
    onChangeModalMode,
    onEditModalClose,
    editModalMode,
  } = useContractModalStateContext();
  const { data } = useGetContractQuery(
    client,
    {
      id: contractId,
    },
    {
      enabled: isEditModalOpen && !!contractId,
      refetchOnMount: true,
    },
  );
  const { data: bankAccountsData } = useBankAccountsQuery(client);
  const { data: tenantSettingsData } = useTenantSettingsQuery(client);

  const queryKey = useGetContractsQuery.getKey({ id: organizationId });
  const contractQueryKey = useGetContractQuery.getKey({ id: organizationId });

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
        queryClient.invalidateQueries({ queryKey: contractQueryKey });
      }, 500);
    },
  });

  const defaultValues = new ContractDetailsDto(data?.contract);

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (_, _action, next) => {
      return next;
    },
  });

  useEffect(() => {
    if (isEditModalOpen) {
      setInitialOpen(editModalMode);
    } else {
      setInitialOpen(EditModalMode.ContractDetails);
    }
  }, [isEditModalOpen]);

  const addressDetailsDefailtValues = new BillingDetailsDto(data?.contract);

  const { state: addressState, setDefaultValues: setDefaultAddressValues } =
    useForm<BillingAddressDetailsFormDto>({
      formId: 'billing-details-address-form',
      defaultValues: addressDetailsDefailtValues,
      stateReducer: (_, _action, next) => {
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
    onEditModalClose();
    onChangeModalMode(EditModalMode.ContractDetails);
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
          onChangeModalMode(EditModalMode.ContractDetails);
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
        },
      ] as unknown as InvoiceLine[],
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
  const saveButtonText = useMemo(() => {
    if (data?.contract?.contractStatus === ContractStatus.Draft) {
      return 'Save draft';
    }

    return 'Save changes';
  }, [data?.contract?.contractStatus]);

  return (
    <Modal open={isEditModalOpen} onOpenChange={handleCloseModal}>
      <ModalPortal>
        <ModalOverlay />
        <ModalContent
          placement='center'
          className='border-r-2 flex gap-6 bg-transparent shadow-none border-none'
          style={{ minWidth: '971px', minHeight: '80vh', boxShadow: 'none' }}
        >
          <div className='relative'>
            <motion.div
              layout
              data-isOpen={editModalMode === EditModalMode.ContractDetails}
              variants={mainVariants as Variants}
              animate={
                editModalMode === EditModalMode.ContractDetails
                  ? 'open'
                  : 'closed'
              }
              onClick={() =>
                editModalMode === EditModalMode.BillingDetails
                  ? onChangeModalMode(EditModalMode.ContractDetails)
                  : null
              }
              className={cn(
                'flex flex-col gap-4 px-6 pb-6 pt-4 bg-white  rounded-lg justify-between relative h-[80vh] min-w-[424px]',
                {
                  'cursor-pointer':
                    editModalMode === EditModalMode.BillingDetails,
                },
              )}
            >
              <ModalHeader className='p-0 font-semibold'>
                <FormInput
                  className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
                  name='contractName'
                  formId={formId}
                />
              </ModalHeader>

              <ContractBillingDetailsForm
                formId={formId}
                contractId={contractId}
                tenantBillingProfile={
                  tenantBillingProfile
                    ?.tenantBillingProfiles?.[0] as TenantBillingProfile
                }
                renewedAt={renewsAt}
                currency={state?.values?.currency?.value}
                bankAccounts={bankAccountsData?.bankAccounts as BankAccount[]}
                payAutomatically={state?.values?.payAutomatically}
                billingEnabled={
                  tenantSettingsData?.tenantSettings?.billingEnabled
                }
              />
              <ModalFooter className='p-0 flex'>
                <Button
                  variant='outline'
                  colorScheme='gray'
                  onClick={handleCloseModal}
                  className='w-full'
                  size='md'
                >
                  Cancel changes
                </Button>
                <Button
                  className='ml-3 w-full'
                  size='md'
                  variant='outline'
                  colorScheme='primary'
                  onClick={handleApplyChanges}
                  loadingText='Saving...'
                  isLoading={updateContract.isPending}
                >
                  {saveButtonText}
                </Button>
              </ModalFooter>
            </motion.div>
            <motion.div
              layout
              data-isOpen={editModalMode === EditModalMode.BillingDetails}
              variants={variants}
              animate={
                editModalMode === EditModalMode.BillingDetails
                  ? 'open'
                  : 'closed'
              }
              className='flex flex-col gap-4 px-6 pb-6 pt-4 bg-white rounded-lg justify-between relative shadow-2xl h-full min-w-[424px]'
            >
              <motion.div
                className='h-full flex flex-col justify-between'
                animate={{
                  opacity:
                    editModalMode === EditModalMode.BillingDetails ? 1 : 0,
                  transition: { duration: 0.2 },
                }}
              >
                <div className='flex flex-col relative justify-between'>
                  <ModalHeader className='p-0 text-lg font-semibold'>
                    <div>
                      {data?.contract?.organizationLegalName ||
                        organizationName ||
                        "Unnamed's "}{' '}
                    </div>
                    <span className='text-base font-normal'>
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
                    onClick={() =>
                      initialOpen === EditModalMode.BillingDetails
                        ? handleCloseModal()
                        : onChangeModalMode(EditModalMode.ContractDetails)
                    }
                    className='w-full'
                    size='md'
                  >
                    Cancel changes
                  </Button>
                  <Button
                    className='ml-3 w-full'
                    size='md'
                    variant='outline'
                    colorScheme='primary'
                    loadingText='Saving...'
                    isLoading={updateContract.isPending}
                    onClick={handleSaveAddressChanges}
                  >
                    {saveButtonText}
                  </Button>
                </ModalFooter>
              </motion.div>
            </motion.div>
          </div>

          {tenantSettingsData?.tenantSettings?.billingEnabled && (
            <div style={{ minWidth: '600px' }} className='bg-white rounded'>
              <div className='w-full h-full'>
                <Invoice
                  shouldBlurDummy
                  onOpenAddressDetailsModal={() =>
                    onChangeModalMode(EditModalMode.BillingDetails)
                  }
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
                  check={
                    tenantBillingProfile?.tenantBillingProfiles?.[0]?.check
                  }
                  availableBankAccount={
                    bankAccountsData?.bankAccounts?.find(
                      (e) => e.currency === state?.values?.currency?.value,
                    ) as BankAccount
                  }
                />
              </div>
            </div>
          )}
        </ModalContent>
      </ModalPortal>
    </Modal>
  );
};
