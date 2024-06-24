import { useParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect } from 'react';

import { motion, Variants } from 'framer-motion';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ModalFooter, ModalHeader } from '@ui/overlay/Modal/Modal';
import { calculateMaxArr } from '@organization/components/Tabs/panels/AccountPanel/utils.ts';
import {
  Contract,
  ContractStatus,
  ServiceLineItem,
  TenantBillingProfile,
} from '@graphql/types';
import { BillingDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm';
import { ContractStatusTag } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ContractStatusTag.tsx';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';
import { ModalWithInvoicePreview } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ModalWithInvoicePreview';

import { ContractBillingDetailsForm } from './ContractBillingDetailsForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  notes?: string | null;
  status: ContractStatus;
  opportunityId?: string;
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
  status,
  serviceStarted,
  opportunityId,
}: SubscriptionServiceModalProps) => {
  const store = useStore();
  const organizationId = useParams().id as string;
  const contractStore = store.contracts.value.get(contractId);
  const contractNameInputRef = useRef<HTMLInputElement | null>(null);
  const organizationStore = store.organizations.value.get(organizationId);
  const opportunityStore = opportunityId
    ? store.opportunities.value.get(opportunityId)
    : undefined;

  const [initialOpen, setInitialOpen] = useState(EditModalMode.ContractDetails);
  const [historyAddressDetails, setHistoryAddressDetails] = useState(
    contractStore?.value?.billingDetails,
  );
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

  useEffect(() => {
    if (isEditModalOpen) {
      setHistoryAddressDetails(contractStore?.value?.billingDetails);
    }
  }, [isEditModalOpen, editModalMode]);
  const getOpportunitiesStores = (contract: Contract) => {
    const c = store.contracts.value.get(contract.metadata.id)?.value;

    return (
      c?.opportunities?.map((e) => {
        return store.opportunities.value.get(e?.metadata?.id)?.value;
      }) || []
    );
  };
  const handleUpdateArrForecast = () => {
    // update opportunity
    const contractLineItemsStores =
      contractStore?.value?.contractLineItems?.map((e) => {
        return contractLineItemsStore.value.get(e.metadata.id)?.value;
      });
    const arrOpportunity = calculateMaxArr(
      contractLineItemsStores as ServiceLineItem[],
      contractStore?.value as Contract,
    );

    opportunityStore?.update(
      (prev) => ({
        ...prev,
        maxAmount: arrOpportunity,
        amount: (arrOpportunity * prev.renewalAdjustedRate) / 100,
      }),
      { mutate: false },
    );

    const organization = organizationStore?.value;
    const contracts = organization?.contracts || [];

    const totalArr = contracts.reduce(
      (acc, contract) => {
        const opportunities = getOpportunitiesStores(contract).filter(
          (e) => e?.internalStage === 'OPEN' && e?.internalType === 'RENEWAL',
        );

        const amount = opportunities.reduce(
          (acc, opportunity) => acc + (opportunity?.amount ?? 0) || 0,
          0,
        );
        const maxAmount = opportunities.reduce(
          (acc, opportunity) => acc + (opportunity?.maxAmount ?? 0) || 0,
          0,
        );

        return {
          maxArrForecast: acc.maxArrForecast + maxAmount,
          arrForecast: acc.arrForecast + amount,
        };
      },
      { maxArrForecast: 0, arrForecast: 0 },
    );

    organizationStore?.update(
      (prev) => ({
        ...prev,
        accountDetails: {
          ...prev.accountDetails,
          renewalSummary: {
            ...prev?.accountDetails?.renewalSummary,
            arrForecast: totalArr.arrForecast,
            maxArrForecast: totalArr.maxArrForecast,
          },
        },
      }),
      { mutate: false },
    );
  };
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
      if (itemStore.value.closed) {
        contractLineItemsStore.closeServiceLineItem(
          { id: e.metadata.id },
          contractId,
        );

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
    handleUpdateArrForecast();
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
                onClick={() => {
                  if (initialOpen === EditModalMode.BillingDetails) {
                    handleCloseModal();

                    return;
                  }
                  onChangeModalMode(EditModalMode.ContractDetails);
                  contractStore?.update((prev) => ({
                    ...prev,
                    billingDetails: historyAddressDetails,
                  }));
                }}
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
