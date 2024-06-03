import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@spaces/utils/date';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { Modal, ModalContent, ModalOverlay } from '@ui/overlay/Modal/Modal';
import { ContractStartModal } from '@organization/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals/ContractStartModal';
import { ContractRenewsModal } from '@organization/components/Tabs/panels/AccountPanel/Contract/ChangeContractStatusModals/ContractRenewModal';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface SubscriptionServiceModalProps {
  contractId: string;
  serviceStarted?: string;
  organizationName: string;
}

export const ContractStatusModal = ({
  contractId,
  organizationName,
  serviceStarted,
}: SubscriptionServiceModalProps) => {
  const { isModalOpen, onStatusModalClose, mode, nextInvoice } =
    useContractModalStatusContext();

  return (
    <Modal
      open={isModalOpen && mode !== ContractStatusModalMode.End}
      onOpenChange={onStatusModalClose}
    >
      <ModalOverlay className='z-50' />
      <ModalContent
        placement={nextInvoice ? 'center' : 'top'}
        className='border-r-2 flex gap-6 bg-transparent shadow-none border-none z-[999]'
        style={{
          minWidth: nextInvoice ? '971px' : 'auto',
          minHeight: nextInvoice ? '80vh' : 'auto',
          boxShadow: 'none',
        }}
      >
        <div
          className={cn(
            'flex flex-col gap-4 px-6 pb-6 pt-4 bg-white  rounded-lg justify-between relative h-full min-w-[424px]',
            {
              'h-[80vh]': nextInvoice,
            },
          )}
        >
          {mode === ContractStatusModalMode.Start && (
            <ContractStartModal
              onClose={onStatusModalClose}
              contractId={contractId}
              organizationName={organizationName}
              serviceStarted={serviceStarted}
            />
          )}

          {mode === ContractStatusModalMode.Renew && (
            <ContractRenewsModal
              onClose={onStatusModalClose}
              contractId={contractId}
            />
          )}
        </div>

        {nextInvoice && (
          <div
            style={{ minWidth: '600px' }}
            className='bg-white rounded relative'
          >
            <p className='absolute top-[-30px] left-3 text-white text-base'>
              <span className='font-semibold mr-1'>Monthly recurring •</span>
              {DateTimeUtils.format(
                nextInvoice.issued,
                DateTimeUtils.dateWithAbreviatedMonth,
              )}
            </p>
            <div className='w-full h-full'>
              <Invoice
                note={nextInvoice?.note}
                invoiceNumber={nextInvoice?.invoiceNumber}
                currency={nextInvoice?.currency}
                billedTo={{
                  addressLine1: nextInvoice.customer.addressLine1 ?? '',
                  addressLine2: nextInvoice.customer.addressLine2 ?? '',
                  locality: nextInvoice.customer.addressLocality ?? '',
                  zip: nextInvoice.customer.addressZip ?? '',
                  country: nextInvoice.customer.addressCountry ?? '',
                  email: nextInvoice.customer.email ?? '',
                  name: nextInvoice.customer.name ?? '',
                  region: nextInvoice.customer.addressRegion ?? '',
                }}
                from={{
                  addressLine1: nextInvoice.provider.addressLine1 ?? '',
                  addressLine2: nextInvoice.provider.addressLine2 ?? '',
                  locality: nextInvoice.provider.addressLocality ?? '',
                  zip: nextInvoice.provider.addressZip ?? '',
                  country: nextInvoice.provider.addressCountry ?? '',
                  email: '',
                  name: nextInvoice.provider.name ?? '',
                  region: nextInvoice.provider.addressRegion ?? '',
                }}
                invoicePeriodStart={nextInvoice?.invoicePeriodStart}
                invoicePeriodEnd={nextInvoice?.invoicePeriodEnd}
                tax={nextInvoice.taxDue}
                lines={nextInvoice?.invoiceLineItems ?? []}
                subtotal={nextInvoice?.subtotal}
                issueDate={nextInvoice?.issued}
                dueDate={nextInvoice?.due}
                total={nextInvoice?.amountDue}
                canPayWithBankTransfer={true}
                check={true}
                availableBankAccount={null}
              />
            </div>
          </div>
        )}
      </ModalContent>
    </Modal>
  );
};
