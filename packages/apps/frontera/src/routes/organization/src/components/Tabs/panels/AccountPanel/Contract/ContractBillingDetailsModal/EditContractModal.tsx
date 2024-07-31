import { useRef, useEffect } from 'react';

import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { useContractModalStateContext } from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

import { BillingAddressModal } from './components/BillingAddressDetails';
import { ContractDetailsModal } from './components/ContractDetails/ContractDetailsModal.tsx';

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

export const EditContractModal = ({
  contractId,
  organizationName,
  status,
  serviceStarted,
  opportunityId,
}: SubscriptionServiceModalProps) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId) as ContractStore;
  const contractNameInputRef = useRef<HTMLInputElement | null>(null);

  const { isEditModalOpen, onEditModalClose } = useContractModalStateContext();

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
            <ContractDetailsModal
              status={status}
              contractId={contractId}
              opportunityId={opportunityId}
              serviceStarted={serviceStarted}
            />
            <BillingAddressModal
              contractId={contractId}
              organizationName={organizationName}
            />
          </div>
        </ModalContent>
      </ModalPortal>
    </Modal>
  );
};
