import { useMemo, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Contract, ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';
import { Card, CardFooter, CardHeader } from '@ui/presentation/Card/Card';
import { UpcomingInvoices } from '@organization/components/Tabs/panels/AccountPanel/Contract/UpcomingInvoices/UpcomingInvoices';
import { useUpdatePanelModalStateContext } from '@organization/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import {
  EditModalMode,
  useContractModalStateContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractModalsContext';

import { Services } from './Services/Services';
import { ContractSubtitle } from './ContractSubtitle';
import { ContractCardActions } from './ContractCardActions';
import { RenewalARRCard } from './RenewalARR/RenewalARRCard';
import { EditContractModal } from './ContractBillingDetailsModal/EditContractModal';

interface ContractCardProps {
  values: Contract;
  organizationName: string;
}

export const ContractCard = observer(
  ({ organizationName, values }: ContractCardProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(values.metadata.id);

    const [isExpanded, setIsExpanded] = useState(
      !contractStore?.value?.contractSigned,
    );
    const { setIsPanelModalOpen } = useUpdatePanelModalStateContext();
    const {
      isEditModalOpen,
      onEditModalOpen,
      onChangeModalMode,
      onEditModalClose,
    } = useContractModalStateContext();

    // this is needed to block scroll on safari when modal is open, scrollbar overflow issue
    useEffect(() => {
      if (isEditModalOpen) {
        setIsPanelModalOpen(true);
      }

      if (!isEditModalOpen) {
        setIsPanelModalOpen(false);
      }
    }, [isEditModalOpen]);

    if (!contractStore) return null;

    const contract = contractStore.value;

    const handleOpenBillingDetails = () => {
      onChangeModalMode(EditModalMode.BillingDetails);
      onEditModalOpen();
    };

    const handleOpenContractDetails = () => {
      onChangeModalMode(EditModalMode.ContractDetails);
      onEditModalOpen();
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
          role='button'
          className='p-0 w-full flex flex-col'
          onClick={() => (!isExpanded ? setIsExpanded(true) : null)}
        >
          <article className='flex justify-between flex-1 w-full'>
            <Input
              name='contractName'
              value={contract?.contractName}
              placeholder='Add contract name'
              onFocus={(e) => e.target.select()}
              className='font-semibold hover:border-none focus:border-none max-h-6 min-h-0 w-full overflow-hidden overflow-ellipsis border-0'
              onChange={(e) =>
                contractStore?.update((prev) => ({
                  ...prev,
                  contractName: e.target.value,
                }))
              }
            />

            <ContractCardActions
              status={contract?.contractStatus}
              contractId={contract?.metadata?.id}
              serviceStarted={contract?.serviceStarted}
              onOpenEditModal={handleOpenContractDetails}
              organizationName={
                contract?.billingDetails?.organizationLegalName ||
                organizationName ||
                'Unnamed'
              }
            />
          </article>

          <div
            tabIndex={1}
            role='button'
            className='w-full'
            onClick={handleOpenContractDetails}
          >
            <ContractSubtitle id={contract.metadata.id} />
          </div>
        </CardHeader>

        <CardFooter className='p-0 mt-0 w-full flex flex-col'>
          {opportunityId &&
            !!contract?.contractLineItems?.filter(
              (e) => !e.metadata.id.includes('new'),
            )?.length && (
              <RenewalARRCard
                currency={contract?.currency}
                opportunityId={opportunityId}
                contractId={contract?.metadata?.id}
                startedAt={contract?.serviceStarted}
                hasEnded={contract?.contractStatus === ContractStatus.Ended}
              />
            )}
          <Services
            id={contract?.metadata?.id}
            currency={contract?.currency}
            onModalOpen={onEditModalOpen}
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
            isOpen={isEditModalOpen}
            onClose={onEditModalClose}
            opportunityId={opportunityId}
            status={contract?.contractStatus}
            contractId={contract?.metadata?.id}
            organizationName={organizationName}
            serviceStarted={contract?.serviceStarted}
            notes={contract?.billingDetails?.invoiceNote}
          />
        </CardFooter>
      </Card>
    );
  },
);
