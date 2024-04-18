import React, { useMemo } from 'react';

import { UseMutationResult } from '@tanstack/react-query';

import { DotLive } from '@ui/media/icons/DotLive';
import { XSquare } from '@ui/media/icons/XSquare';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import { ContractEndModal } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ChangeContractStatusModals';
import { ContractMenu } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCardActions/ContractMenu';
import { ContractStatusTag } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCardActions/ContractStatusTag';
import { ContractStatusModal } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ChangeContractStatusModals/ContractStatusModal';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractStatusSelectProps {
  renewsAt?: string;
  contractId: string;
  status: ContractStatus;
  serviceStarted?: string;
  upcomingInvoices: any[];
  organizationName: string;
  nextInvoiceDate?: string;
  onOpenEditModal: () => void;

  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

export const ContractCardActions: React.FC<ContractStatusSelectProps> = ({
  status,
  renewsAt,
  contractId,
  organizationName,
  serviceStarted,
  onUpdateContract,
  nextInvoiceDate,
  onOpenEditModal,
  upcomingInvoices,
}) => {
  const { onStatusModalOpen } = useContractModalStatusContext();

  const getStatusDisplay = useMemo(() => {
    let icon, text;
    switch (status) {
      case ContractStatus.Live:
        icon = <XSquare color='gray.500' mr={1} />;
        text = 'End contract...';
        break;
      case ContractStatus.Draft:
      case ContractStatus.Ended:
        icon = <DotLive color='gray.500' mr={1} />;
        text = 'Make live';
        break;
      case ContractStatus.OutOfContract:
        icon = <RefreshCw05 color='gray.500' mr={2} />;
        text = 'Renew contract';
        break;
      case ContractStatus.Scheduled:
        icon = null;
        text = null;
        break;
      default:
        icon = null;
        text = null;
    }

    return (
      <>
        {icon}
        {text}
      </>
    );
  }, [status]);
  console.log('🏷️ ----- upcomingInvoices: ', upcomingInvoices);
  const handleChangeStatus = () => {
    switch (status) {
      case ContractStatus.Live:
        onStatusModalOpen(ContractStatusModalMode.End);
        break;
      case ContractStatus.Draft:
      case ContractStatus.Ended:
        onStatusModalOpen(ContractStatusModalMode.Start);

        break;
      case ContractStatus.OutOfContract:
        onStatusModalOpen(ContractStatusModalMode.Renew);

        break;
      case ContractStatus.Scheduled:
        break;
      default:
    }
  };

  return (
    <div className='flex items-center gap-2 ml-2'>
      <ContractStatusTag
        status={status}
        contractStarted={serviceStarted}
        statusContent={getStatusDisplay}
        onHandleStatusChange={handleChangeStatus}
      />
      <ContractMenu
        onOpenEditModal={onOpenEditModal}
        statusContent={getStatusDisplay}
        onHandleStatusChange={handleChangeStatus}
        status={status}
      />
      <ContractEndModal
        contractId={contractId}
        organizationName={organizationName}
        renewsAt={renewsAt}
        serviceStarted={serviceStarted}
        nextInvoiceDate={nextInvoiceDate}
        onUpdateContract={onUpdateContract}
      />
      <ContractStatusModal
        contractId={contractId}
        organizationName={organizationName}
        renewsAt={renewsAt}
        serviceStarted={serviceStarted}
        onUpdateContract={onUpdateContract}
        nextInvoiceDate={nextInvoiceDate}
        upcomingInvoices={upcomingInvoices}
      />
    </div>
  );
};
