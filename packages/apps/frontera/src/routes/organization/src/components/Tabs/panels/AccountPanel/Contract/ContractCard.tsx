import { useMemo, useState, useEffect } from 'react';

import { reaction, comparer } from 'mobx';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider/Divider';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Card, CardFooter, CardHeader } from '@ui/presentation/Card/Card';
import { UpcomingInvoices } from '@organization/components/Tabs/panels/AccountPanel/Contract/UpcomingInvoices/UpcomingInvoices';
// import { useUpdatePanelModalStateContext } from '@organization/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import {
  EditModalMode,
  // useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

import { Services } from './Services/Services';
import { ContractSubtitle } from './ContractSubtitle';
import { ContractCardActions } from './ContractCardActions';
import { RenewalARRCard } from './RenewalARR/RenewalARRCard';
import { EditContractModal } from './ContractBillingDetailsModal/EditContractModal';

interface ContractCardProps {
  id: string;
  organizationName: string;
}

export const ContractCard = observer(
  ({ organizationName, id }: ContractCardProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(id);
    const contractLineItemsStore = store.contractLineItems;

    const { onClose, onOpen, open } = useDisclosure({
      id: `contract-card-modal-${id}`,
    });
    const [mode, setMode] = useState<EditModalMode>(
      EditModalMode.ContractDetails,
    );
    // const [isPanelModalOpen, setIsPanelModalOpen] = useState(false);

    const [isExpanded, setIsExpanded] = useState(
      !contractStore?.value?.contractSigned,
    );
    // const { setIsPanelModalOpen } = useUpdatePanelModalStateContext();
    // const {
    //   isEditModalOpen,
    //   onEditModalOpen,
    //   onChangeModalMode,
    //   onEditModalClose,
    // } = useContractModalStateContext();

    // this is needed to block scroll on safari when modal is open, scrollbar overflow issue
    // useEffect(() => {
    //   if (open) {
    //     setIsPanelModalOpen(true);
    //   } else {
    //     setIsPanelModalOpen(false);
    //   }
    // }, [open]);

    useEffect(() => {
      const dispose = reaction(
        () => contractLineItemsStore.value,
        () => {
          // simulate invoices
          // console.log('🏷️ ----- : SIMULATING INVOICES');
        },
        { equals: comparer.structural },
      );

      return () => dispose();
    }, []);

    if (!contractStore) return null;

    const contract = contractStore.value;

    const handleOpenBillingDetails = () => {
      setMode(EditModalMode.BillingDetails);
      onOpen();
    };
    const handleOpenContractDetails = () => {
      setMode(EditModalMode.ContractDetails);
      onOpen();
    };

    const opportunityId = useMemo(() => {
      return (
        contract?.opportunities?.find((e) => e.internalStage === 'OPEN')?.id ||
        contract?.opportunities?.[0]?.id
      );
    }, []);

    return (
      <Card className='px-4 py-3 w-full text-lg bg-gray-50 transition-all-0.2s-ease-out border border-gray-200 text-gray-700 '>
        <CardHeader
          className='p-0 w-full flex flex-col'
          role='button'
          onClick={() => (!isExpanded ? setIsExpanded(true) : null)}
        >
          <article className='flex justify-between flex-1 w-full'>
            <Input
              className='font-semibold no-border-bottom hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis'
              name='contractName'
              placeholder='Add contract name'
              value={contract?.contractName}
              onChange={(e) =>
                contractStore?.update((value) => {
                  value.contractName = e.target.value;

                  return value;
                })
              }
              onFocus={(e) => e.target.select()}
            />

            <ContractCardActions
              onOpenEditModal={handleOpenContractDetails}
              status={contract?.contractStatus}
              contractId={contract?.metadata?.id}
              serviceStarted={contract?.serviceStarted}
              organizationName={
                contract?.billingDetails?.organizationLegalName ||
                organizationName ||
                'Unnamed'
              }
            />
          </article>

          <div
            role='button'
            tabIndex={1}
            onClick={handleOpenContractDetails}
            className='w-full'
          >
            <ContractSubtitle data={contract} />
          </div>
        </CardHeader>

        <CardFooter className='p-0 mt-0 w-full flex flex-col'>
          {opportunityId && !!contract?.contractLineItems?.length && (
            <RenewalARRCard
              contractId={contract?.metadata?.id}
              hasEnded={contract?.contractStatus === ContractStatus.Ended}
              startedAt={contract?.serviceStarted}
              currency={contract?.currency}
              opportunityId={opportunityId}
            />
          )}
          <Services
            onModalOpen={onOpen}
            currency={contract?.currency}
            data={contract?.contractLineItems}
          />
          {!!contract?.upcomingInvoices?.length && (
            <>
              <Divider className='my-3' />
              <UpcomingInvoices
                data={contract}
                onOpenBillingDetailsModal={handleOpenBillingDetails}
                onOpenServiceLineItemsModal={handleOpenContractDetails}
              />
            </>
          )}

          <EditContractModal
            mode={mode}
            isOpen={open}
            onClose={onClose}
            onChangeMode={setMode}
            status={contract?.contractStatus}
            contractId={contract?.metadata?.id}
            serviceStarted={contract?.serviceStarted}
            organizationName={organizationName}
            notes={contract?.billingDetails?.invoiceNote}
            renewsAt={contract?.opportunities?.[0]?.renewedAt}
          />
        </CardFooter>
      </Card>
    );
  },
);
