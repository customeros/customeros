import { useParams } from 'react-router-dom';
import { useRef, useMemo, useEffect } from 'react';

import { motion, Variants } from 'framer-motion';
import { ContractStore } from '@store/Contracts/Contract.store.ts';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { calculateMaxArr } from '@organization/components/Tabs/panels/AccountPanel/utils.ts';
import {
  Contract,
  ContractStatus,
  ServiceLineItem,
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
import { BillingDetailsForm } from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm';
import { ContractStatusTag } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ContractStatusTag.tsx';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

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
  const contractStore = store.contracts.value.get(contractId) as ContractStore;
  const contractNameInputRef = useRef<HTMLInputElement | null>(null);
  const organizationStore = store.organizations.value.get(organizationId);
  const opportunityStore = opportunityId
    ? store.opportunities.value.get(opportunityId)
    : undefined;

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
      setTimeout(() => {
        contractNameInputRef.current?.focus();
        contractNameInputRef.current?.select();
      });
    }
  }, [isEditModalOpen]);

  useEffect(() => {
    if (isEditModalOpen) {
      contractStore?.setTempValue();
    }
  }, [isEditModalOpen]);

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
      contractStore?.tempValue as Contract,
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
  const handleApplyChanges = async () => {
    contractStore.value = contractStore.tempValue;
    await contractStore.updateContractValues();

    contractStore?.value?.contractLineItems?.forEach((e) => {
      const itemStore = contractLineItemsStore.value.get(
        e.metadata.id,
      ) as ContractLineItemStore;
      itemStore.value = itemStore.tempValue;

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
        contractLineItemsStore.createNewVersion(itemStore?.value, contractId);

        return;
      }

      itemStore?.updateServiceLineItem();
    });
    handleCloseModal();
    handleUpdateArrForecast();
  };

  const handleSaveAddressChanges = async () => {
    contractStore.update(
      (contract) => {
        contract.billingDetails = {
          ...contract.billingDetails,
          organizationLegalName:
            contractStore?.tempValue?.billingDetails?.organizationLegalName,
          country: contractStore?.tempValue?.billingDetails?.country,
          addressLine1: contractStore?.tempValue?.billingDetails?.addressLine1,
          addressLine2: contractStore?.tempValue?.billingDetails?.addressLine2,
          locality: contractStore?.tempValue?.billingDetails?.locality,
          postalCode: contractStore?.tempValue?.billingDetails?.postalCode,
          region: contractStore?.tempValue?.billingDetails?.region,

          billingEmail: contractStore?.tempValue?.billingDetails?.billingEmail,
          billingEmailCC:
            contractStore?.tempValue?.billingDetails?.billingEmailCC,
          billingEmailBCC:
            contractStore?.tempValue?.billingDetails?.billingEmailBCC,
        };

        return contract;
      },
      { mutate: false },
    );

    await contractStore.updateBillingAddress();
    onChangeModalMode(EditModalMode.ContractDetails);
  };

  const handleCancelAddressChanges = () => {
    contractStore?.updateTemp((contract) => {
      contract.billingDetails = {
        ...contract.billingDetails,
        organizationLegalName:
          contractStore?.value?.billingDetails?.organizationLegalName,
        country: contractStore?.value?.billingDetails?.country,
        addressLine1: contractStore?.value?.billingDetails?.addressLine1,
        addressLine2: contractStore?.value?.billingDetails?.addressLine2,
        locality: contractStore?.value?.billingDetails?.locality,
        postalCode: contractStore?.value?.billingDetails?.postalCode,
        region: contractStore?.value?.billingDetails?.region,

        billingEmail: contractStore?.value?.billingDetails?.billingEmail,
        billingEmailCC: contractStore?.value?.billingDetails?.billingEmailCC,
        billingEmailBCC: contractStore?.value?.billingDetails?.billingEmailBCC,
      };

      return contract;
    });

    onChangeModalMode(EditModalMode.ContractDetails);
  };

  const availableCurrencies = useMemo(
    () => bankAccounts?.map((e) => e.value.currency),
    [],
  );

  const canAllowPayWithBankTransfer = useMemo(() => {
    return availableCurrencies.includes(contractStore?.tempValue?.currency);
  }, [availableCurrencies, contractStore?.tempValue?.currency]);

  useEffect(() => {
    if (!canAllowPayWithBankTransfer) {
      contractStore?.updateTemp((prev) => ({
        ...prev,
        billingDetails: {
          ...prev.billingDetails,
          canPayWithBankTransfer: false,
        },
      }));
    }
  }, [canAllowPayWithBankTransfer]);

  return (
    <Modal open={isEditModalOpen} onOpenChange={onEditModalClose}>
      <ModalPortal>
        <ModalOverlay className='z-50' />
        <ModalContent
          placement='center'
          className='border-r-2 flex bg-transparent shadow-none border-none z-[999] w-full '
          style={{
            minWidth: 'auto',
            minHeight: '80vh',
            boxShadow: 'none',
          }}
        >
          <div className='relative '>
            <motion.div
              layout
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
                'flex flex-col gap-4 px-6 pb-6 pt-4 bg-gray-25  rounded-lg justify-between relative h-[80vh] min-w-[460px] overflow-y-auto overflow-x-hidden',
                {
                  'cursor-pointer':
                    editModalMode === EditModalMode.BillingDetails,
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
                  value={contractStore?.tempValue?.contractName}
                  onChange={(e) =>
                    contractStore?.updateTemp((prev) => ({
                      ...prev,
                      contractName: e.target.value,
                    }))
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
                contractStatus={contractStore?.tempValue?.contractStatus}
                openAddressModal={() =>
                  onChangeModalMode(EditModalMode.BillingDetails)
                }
              />
              <ModalFooter className='p-0 flex sticky z-[999] -bottom-6 -mb-5 pb-5 pt-3 bg-gray-25'>
                <Button
                  variant='outline'
                  colorScheme='gray'
                  onClick={() => {
                    handleCloseModal();
                    contractStore?.resetTempValue();
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
                  onClick={() => handleApplyChanges()}
                  loadingText='Saving...'
                >
                  {contractStore?.value?.contractStatus === ContractStatus.Draft
                    ? 'Save draft'
                    : 'Save'}
                </Button>
              </ModalFooter>
            </motion.div>
            <motion.div
              layout
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
                      {contractStore?.tempValue?.billingDetails
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
                      handleCancelAddressChanges();
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
                    Save
                  </Button>
                </ModalFooter>
              </motion.div>
            </motion.div>
          </div>
        </ModalContent>
      </ModalPortal>
    </Modal>
  );
};
