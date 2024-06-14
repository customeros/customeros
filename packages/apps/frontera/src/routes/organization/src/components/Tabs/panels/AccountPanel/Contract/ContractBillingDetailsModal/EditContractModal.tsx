import { useRef, useMemo, useEffect } from 'react';

import { Input } from '@ui/form/Input';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { ContractStatus, TenantBillingProfile } from '@graphql/types';
import { EditModalMode } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';
import { BillingAddressDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm';
import { ModalWithInvoicePreview } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ModalWithInvoicePreview';

import { contractOptionIcon } from '../ContractCardActions/utils';
import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  mode: EditModalMode;
  notes?: string | null;
  status: ContractStatus;
  serviceStarted?: string;
  organizationName: string;
  onChangeMode: (mode: EditModalMode) => void;
}

export const EditContractModal = ({
  mode,
  isOpen,
  status,
  onClose,
  renewsAt,
  contractId,
  onChangeMode,
  serviceStarted,
  organizationName,
}: SubscriptionServiceModalProps) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId);
  const contractNameInputRef = useRef<HTMLInputElement | null>(null);

  const bankAccounts = store.settings.bankAccounts.toArray();
  const tenantSettings = store.settings.tenant.value;
  const tenantBillingProfiles = store.settings.tenantBillingProfiles.toArray();
  const contractLineItemsStore = store.contractLineItems;

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        contractNameInputRef.current?.focus();
        contractNameInputRef.current?.select();
      });
    }
  }, [isOpen]);

  const handleCloseModal = () => {
    onClose();
    onChangeMode(EditModalMode.BillingDetails);
  };
  const handleApplyChanges = () => {
    contractStore?.update((prev) => prev);
    contractStore?.value?.contractLineItems?.forEach((e) => {
      if (e.metadata.id.includes('new') && !e.parentId) {
        contractLineItemsStore.createNewServiceLineItem(e, contractId);

        return;
      }
      if (e.metadata.id.includes('new') && !!e.parentId) {
        contractLineItemsStore.createNewVersion(e);

        return;
      }
      const contractLineItem = contractLineItemsStore.value.get(e.metadata.id);
      contractLineItem?.update((prev) => prev);
    });

    // TODO mutate SLIs that were change during that session
    handleCloseModal();
  };

  const handleSaveAddressChanges = () => {
    contractStore?.update((prev) => ({
      ...prev,
    }));
    onChangeMode(EditModalMode.ContractDetails);
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
      isOpen={isOpen}
      contractId={contractId}
      onClose={handleCloseModal}
      availableBankAccount={availableBankAccount}
      billingEnabled={tenantSettings?.billingEnabled}
      showNextInvoice={tenantSettings?.billingEnabled}
    >
      <ModalBody className='bg-white pr-0 rounded-lg'>
        <ScrollArea>
          <div className='relative flex flex-col space-y-2 pt-4'>
            {mode === EditModalMode.ContractDetails ? (
              <>
                <div className='flex justify-between'>
                  <Input
                    ref={contractNameInputRef}
                    variant='unstyled'
                    className='font-semibold max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
                    name='contractName'
                    placeholder='Add contract name'
                    onFocus={(e) => e.target.select()}
                    autoComplete='off'
                    value={contractStore?.value?.contractName}
                    onChange={(e) =>
                      contractStore?.update(
                        (values) => {
                          values.contractName = e.target.value;

                          return values;
                        },
                        { mutate: false },
                      )
                    }
                  />

                  <ContractStatusTag
                    status={status}
                    contractStarted={serviceStarted}
                  />
                </div>

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
              </>
            ) : (
              <>
                <div className='flex flex-col w-full justify-between'>
                  <div className='p-0 text-lg font-semibold'>
                    <div>
                      {contractStore?.value?.billingDetails
                        ?.organizationLegalName ||
                        organizationName ||
                        "Unnamed's "}{' '}
                    </div>
                    <span className='text-base font-normal'>
                      These details are required to issue invoices
                    </span>
                  </div>

                  <BillingAddressDetailsForm contractId={contractId} />
                </div>
              </>
            )}

            <div className='fixed bottom-[1px] flex py-6 ring-1 ring-white gap-3 bg-white w-full max-w-[342px]'>
              <Button
                variant='outline'
                colorScheme='gray'
                onClick={
                  mode === EditModalMode.ContractDetails
                    ? handleCloseModal
                    : () => onChangeMode(EditModalMode.ContractDetails)
                }
                className='w-full'
                size='md'
              >
                Cancel changes
              </Button>
              <Button
                className='w-full'
                size='md'
                variant='outline'
                colorScheme='primary'
                onClick={
                  mode === EditModalMode.ContractDetails
                    ? handleApplyChanges
                    : handleSaveAddressChanges
                }
                loadingText='Saving...'
              >
                {saveButtonText}
              </Button>
            </div>
          </div>
        </ScrollArea>
      </ModalBody>
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

const ScrollArea = ({ children }: { children: React.ReactNode }) => {
  return (
    <ScrollAreaRoot className='h-[80vh] pr-8'>
      <ScrollAreaViewport>{children}</ScrollAreaViewport>
      <ScrollAreaScrollbar>
        <ScrollAreaThumb />
      </ScrollAreaScrollbar>
    </ScrollAreaRoot>
  );
};
