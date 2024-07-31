import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Modal, ModalContent, ModalOverlay } from '@ui/overlay/Modal/Modal';
import { ContractStartModal } from '@organization/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals/ContractStartModal';
import { ContractRenewsModal } from '@organization/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals/ContractRenewModal';
import { ContractDeleteModal } from '@organization/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals/ContractDeleteModal.tsx';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface SubscriptionServiceModalProps {
  contractId: string;
  serviceStarted?: string;
  organizationName: string;
}

export const ContractStatusModal = observer(
  ({
    contractId,
    organizationName,
    serviceStarted,
  }: SubscriptionServiceModalProps) => {
    const { isModalOpen, onStatusModalClose, mode } =
      useContractModalStatusContext();

    return (
      <Modal
        onOpenChange={onStatusModalClose}
        open={isModalOpen && mode !== ContractStatusModalMode.End}
      >
        <ModalOverlay className='z-50' />
        <ModalContent
          placement={'top'}
          className='border-r-2 flex gap-6 bg-transparent shadow-none border-none z-[999]'
          style={{
            minWidth: 'auto',
            minHeight: 'auto',
            boxShadow: 'none',
          }}
        >
          <div
            className={cn(
              'flex flex-col gap-4 px-6 pb-6 pt-4 bg-white  rounded-lg justify-between relative h-full min-w-[424px]',
            )}
          >
            {mode === ContractStatusModalMode.Start && (
              <ContractStartModal
                contractId={contractId}
                onClose={onStatusModalClose}
                serviceStarted={serviceStarted}
                organizationName={organizationName}
              />
            )}

            {mode === ContractStatusModalMode.Renew && (
              <ContractRenewsModal
                contractId={contractId}
                onClose={onStatusModalClose}
              />
            )}
            {mode === ContractStatusModalMode.Delete && (
              <ContractDeleteModal
                contractId={contractId}
                onClose={onStatusModalClose}
              />
            )}
          </div>
        </ModalContent>
      </Modal>
    );
  },
);
