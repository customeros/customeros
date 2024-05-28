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
import { DateTimeUtils } from '@spaces/utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { ModalFooter, ModalHeader } from '@ui/overlay/Modal/Modal';
import { useGetContractQuery } from '@organization/graphql/getContract.generated';
import {
  BankAccount,
  ContractStatus,
  TenantBillingProfile,
} from '@graphql/types';
import { useUpdateContractMutation } from '@organization/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import { BillingDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { ModalWithInvoicePreview } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ModalWithInvoicePreview';
import { useEditContractModalStores } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/EditContractModalStores';
import {
  BillingDetailsDto,
  BillingAddressDetailsFormDto,
} from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm.dto';

import { ContractDetailsDto } from './ContractDetails.dto';
import { contractOptionIcon } from '../ContractCardActions/utils';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  status: ContractStatus;
  serviceStarted?: string;
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
  renewsAt,
  status,
  serviceStarted,
}: SubscriptionServiceModalProps) => {
  const formId = `billing-details-form-${contractId}`;
  const organizationId = useParams()?.id as string;
  const client = getGraphQLClient();
  const { serviceFormStore } = useEditContractModalStores();

  const contractNameInputRef = useRef<HTMLInputElement | null>(null);

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
      setTimeout(() => {
        contractNameInputRef.current?.focus();
        contractNameInputRef.current?.select();
      });
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

    const savingServiceLineItems = serviceFormStore.saveServiceLineItems();
    const updatingContract = new Promise((resolve, reject) => {
      updateContract.mutate(
        {
          input: {
            contractId,
            ...payload,
            contractName:
              payload.contractName?.length === 0
                ? 'Unnamed contract'
                : payload.contractName,
          },
        },
        {
          onSuccess: () => resolve('Billing details updated'),
          onError: reject,
        },
      );
    });

    Promise.all([savingServiceLineItems, updatingContract])
      .then(() => {
        toastSuccess(
          'Contract details updated',
          `update-contract-success-${contractId}`,
        );
        handleCloseModal();
      })
      .catch(() => {
        toastError('Update failed', `update-contract-error-${contractId}`);
      });
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

  const availableCurrencies = useMemo(
    () => (bankAccountsData?.bankAccounts ?? []).map((e) => e.currency),
    [],
  );

  const canAllowPayWithBankTransfer = useMemo(() => {
    return availableCurrencies.includes(state.values.currency?.value);
  }, [availableCurrencies, state.values.currency]);

  const availableBankAccount = useMemo(
    () =>
      (bankAccountsData?.bankAccounts ?? []).find(
        (e) => e.currency === state.values.currency?.value,
      ),
    [state.values.currency?.value && bankAccountsData?.bankAccounts.length],
  );

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
    <ModalWithInvoicePreview
      currency={state?.values?.currency?.value}
      allowAutoPayment={state?.values?.payAutomatically}
      allowOnlinePayment={state?.values?.payOnline}
      allowBankTransfer={state?.values?.canPayWithBankTransfer}
      availableBankAccount={availableBankAccount as BankAccount}
      allowCheck={state?.values?.check}
      billingEnabled={tenantSettingsData?.tenantSettings?.billingEnabled}
      showNextInvoice={tenantSettingsData?.tenantSettings?.billingEnabled}
    >
      <div className='relative'>
        <motion.div
          layout
          variants={mainVariants as Variants}
          animate={
            editModalMode === EditModalMode.ContractDetails ? 'open' : 'closed'
          }
          onClick={() =>
            editModalMode === EditModalMode.BillingDetails
              ? onChangeModalMode(EditModalMode.ContractDetails)
              : null
          }
          className={cn(
            'flex flex-col gap-4 px-6 pb-6 pt-4 bg-white  rounded-lg justify-between relative h-[80vh] min-w-[424px] overflow-y-auto overflow-x-hidden',
            {
              'cursor-pointer': editModalMode === EditModalMode.BillingDetails,
            },
          )}
        >
          <ModalHeader className='p-0 font-semibold flex'>
            <FormInput
              ref={contractNameInputRef}
              className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
              name='contractName'
              placeholder='Add contract name'
              formId={formId}
              onFocus={(e) => e.target.select()}
            />

            <ContractStatusTag
              status={status}
              contractStarted={serviceStarted}
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
            billingEnabled={tenantSettingsData?.tenantSettings?.billingEnabled}
            contractStatus={data?.contract?.contractStatus}
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
              onClick={() => handleApplyChanges()}
              loadingText='Saving...'
              isLoading={updateContract.isPending}
            >
              {saveButtonText}
            </Button>
          </ModalFooter>
        </motion.div>
        <motion.div
          layout
          variants={variants}
          animate={
            editModalMode === EditModalMode.BillingDetails ? 'open' : 'closed'
          }
          className='flex flex-col gap-4 px-6 pb-6 pt-4 bg-white rounded-lg justify-between relative shadow-2xl h-full min-w-[424px]'
        >
          <motion.div
            className='h-full flex flex-col justify-between'
            animate={{
              opacity: editModalMode === EditModalMode.BillingDetails ? 1 : 0,
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
    </ModalWithInvoicePreview>
  );
};

const ContractStatusTag = ({
  status,
  contractStarted,
}: {
  status: ContractStatus;
  contractStarted?: string;
}) => {
  const statusColorScheme: Record<string, string> = {
    [ContractStatus.Live]: 'primary',
    [ContractStatus.Draft]: 'gray',
    [ContractStatus.Ended]: 'gray',
    [ContractStatus.Scheduled]: 'primary',
    [ContractStatus.OutOfContract]: 'warning',
  };
  const contractStatusOptions: SelectOption<ContractStatus>[] = [
    { label: 'Draft', value: ContractStatus.Draft },
    { label: 'Ended', value: ContractStatus.Ended },
    { label: 'Live', value: ContractStatus.Live },
    { label: 'Out of contract', value: ContractStatus.OutOfContract },
    {
      label: contractStarted
        ? `Live ${DateTimeUtils.format(
            contractStarted,
            DateTimeUtils.defaultFormatShortString,
          )}`
        : 'Scheduled',
      value: ContractStatus.Scheduled,
    },
  ];
  const icon = contractOptionIcon?.[status];
  const selected = contractStatusOptions.find((e) => e.value === status);

  return (
    <>
      <Tag
        className='flex items-center gap-1 whitespace-nowrap mx-0 px-1'
        colorScheme={statusColorScheme[status] as 'primary'}
      >
        <TagLeftIcon className='m-0'>{icon}</TagLeftIcon>

        <TagLabel>{selected?.label}</TagLabel>
      </Tag>
    </>
  );
};
