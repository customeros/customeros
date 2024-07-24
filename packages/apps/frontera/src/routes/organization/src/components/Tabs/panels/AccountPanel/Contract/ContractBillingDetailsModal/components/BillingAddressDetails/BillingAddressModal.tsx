import { useRef, useMemo, useEffect } from 'react';

import { motion } from 'framer-motion';
import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { ModalFooter, ModalHeader } from '@ui/overlay/Modal/Modal.tsx';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext.tsx';

import { BillingDetailsForm } from './BillingAddressDetailsForm.tsx';

interface BillingAddressModalProps {
  contractId: string;
  organizationName: string;
}

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

export const BillingAddressModal = observer(
  ({ contractId, organizationName }: BillingAddressModalProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(
      contractId,
    ) as ContractStore;
    const contractNameInputRef = useRef<HTMLInputElement | null>(null);

    const { isEditModalOpen, onChangeModalMode, editModalMode } =
      useContractModalStateContext();

    const bankAccounts = store.settings.bankAccounts.toArray();

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

    const handleSaveAddressChanges = async () => {
      contractStore.update(
        (contract) => {
          contract.billingDetails = {
            ...contract.billingDetails,
            organizationLegalName:
              contractStore?.tempValue?.billingDetails?.organizationLegalName,
            country: contractStore?.tempValue?.billingDetails?.country,
            addressLine1:
              contractStore?.tempValue?.billingDetails?.addressLine1,
            addressLine2:
              contractStore?.tempValue?.billingDetails?.addressLine2,
            locality: contractStore?.tempValue?.billingDetails?.locality,
            postalCode: contractStore?.tempValue?.billingDetails?.postalCode,
            region: contractStore?.tempValue?.billingDetails?.region,

            billingEmail:
              contractStore?.tempValue?.billingDetails?.billingEmail,
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
          billingEmailBCC:
            contractStore?.value?.billingDetails?.billingEmailBCC,
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
      <>
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
      </>
    );
  },
);
