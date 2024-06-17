import { useRef, useMemo, useState, useEffect } from 'react';

import { motion, Variants } from 'framer-motion';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { ModalFooter, ModalHeader } from '@ui/overlay/Modal/Modal';
import { ContractStatus, TenantBillingProfile } from '@graphql/types';
import { BillingDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { ModalWithInvoicePreview } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ModalWithInvoicePreview';

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
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId);
  const contractNameInputRef = useRef<HTMLInputElement | null>(null);

  const [initialOpen, setInitialOpen] = useState(EditModalMode.ContractDetails);
  useState<boolean>(false);
  const {
    isEditModalOpen,
    onChangeModalMode,
    onEditModalClose,
    editModalMode,
  } = useContractModalStateContext();

  const bankAccounts = store.settings.bankAccounts.toArray();
  const tenantSettings = store.settings.tenant.value;
  const tenantBillingProfiles = store.settings.tenantBillingProfiles.toArray();
  const contractLineItemsStore = store.contractLineItems;

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

  const handleCloseModal = () => {
    onEditModalClose();
    onChangeModalMode(EditModalMode.ContractDetails);
  };
  const handleApplyChanges = () => {
    contractStore?.update((prev) => prev);
    contractStore?.value?.contractLineItems?.forEach((e) => {
      const itemStore = contractLineItemsStore.value.get(e.metadata.id);
      if (!itemStore?.value) {
        return;
      }
      if (e.metadata.id.includes('new') && !e.parentId) {
        contractLineItemsStore.createNewServiceLineItem(
          itemStore?.value,
          contractId,
        );

        return;
      }
      if (e.metadata.id.includes('new') && !!e.parentId) {
        contractLineItemsStore.createNewVersion(itemStore?.value);

        return;
      }
      itemStore?.update((prev) => prev);
    });

    handleCloseModal();
  };

  const handleSaveAddressChanges = () => {
    contractStore?.update((prev) => ({
      ...prev,
    }));
    onChangeModalMode(EditModalMode.ContractDetails);
  };

  const availableCurrencies = useMemo(
    () => bankAccounts?.map((e) => e.value.currency),
    [],
  );

  const canAllowPayWithBankTransfer = useMemo(() => {
    return availableCurrencies.includes(contractStore?.value?.currency);
  }, [availableCurrencies, contractStore?.value?.currency]);

  const availableBankAccount = useMemo(
    () =>
      (bankAccounts ?? [])?.find(
        (e) => e?.value?.currency === contractStore?.value?.currency,
      ),
    [contractStore?.value?.currency && bankAccounts],
  );

  useEffect(() => {
    if (!canAllowPayWithBankTransfer) {
      contractStore?.update(
        (prev) => ({
          ...prev,
          billingDetails: {
            ...prev.billingDetails,
            canPayWithBankTransfer: false,
          },
        }),
        { mutate: false },
      );
    }
  }, [canAllowPayWithBankTransfer]);
  const saveButtonText = useMemo(() => {
    if (contractStore?.value?.contractStatus === ContractStatus.Draft) {
      return 'Save draft';
    }

    return 'Save changes';
  }, [contractStore?.value?.contractStatus]);

  return (
    <ModalWithInvoicePreview
      contractId={contractId}
      availableBankAccount={availableBankAccount}
      billingEnabled={tenantSettings?.billingEnabled}
      showNextInvoice={tenantSettings?.billingEnabled}
    >
      <div className='relative '>
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
            'flex flex-col gap-4 px-6 pb-6 pt-4 bg-gray-25  rounded-lg justify-between relative h-[80vh] min-w-[424px] overflow-y-auto overflow-x-hidden',
            {
              'cursor-pointer': editModalMode === EditModalMode.BillingDetails,
            },
          )}
        >
          <ModalHeader className='p-0 font-semibold flex'>
            <Input
              ref={contractNameInputRef}
              className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
              name='contractName'
              placeholder='Add contract name'
              onFocus={(e) => e.target.select()}
              value={contractStore?.value?.contractName}
              onChange={(e) =>
                contractStore?.update(
                  (prev) => ({
                    ...prev,
                    contractName: e.target.value,
                  }),
                  { mutate: false },
                )
              }
            />

            <ContractStatusTag
              status={status}
              contractStarted={serviceStarted}
            />
          </ModalHeader>

          <ContractBillingDetailsForm
            contractId={contractId}
            tenantBillingProfile={
              tenantBillingProfiles?.[0]?.value as TenantBillingProfile
            }
            renewedAt={renewsAt}
            bankAccounts={bankAccounts}
            billingEnabled={tenantSettings?.billingEnabled}
            contractStatus={contractStore?.value?.contractStatus}
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
                  {contractStore?.value?.billingDetails
                    ?.organizationLegalName ||
                    organizationName ||
                    "Unnamed's "}{' '}
                </div>
                <span className='text-base font-normal'>
                  These details are required to issue invoices
                </span>
              </ModalHeader>

              <BillingDetailsForm contractId={contractId} />
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
