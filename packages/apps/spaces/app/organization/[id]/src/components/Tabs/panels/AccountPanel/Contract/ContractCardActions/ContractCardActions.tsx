import React, { useMemo } from 'react';

import { UseMutationResult } from '@tanstack/react-query';

import { useDisclosure } from '@ui/utils';
import { DotLive } from '@ui/media/icons/DotLive';
import { XSquare } from '@ui/media/icons/XSquare';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import { useRenewContractMutation } from '@organization/src/graphql/renewContract.generated';
import { ContractMenu } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCardActions/ContractMenu';
import { ContractStatusTag } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractCardActions/ContractStatusTag';
import {
  ContractEndModal,
  ContractStartModal,
} from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ChangeContractStatusModals';

interface ContractStatusSelectProps {
  renewsAt?: string;
  contractId: string;
  status: ContractStatus;
  serviceStarted?: string;
  organizationName: string;
  nextInvoiceDate?: string;
  contractStarted?: string;
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
  contractStarted,
}) => {
  const {
    onOpen: onOpenEndModal,
    onClose: onCloseEndModal,
    isOpen,
  } = useDisclosure({
    id: 'end-contract-modal',
  });
  const {
    onOpen: onOpenStartModal,
    onClose: onCloseStartModal,
    isOpen: isStartModalOpen,
  } = useDisclosure({
    id: 'start-contract-modal',
  });
  const client = getGraphQLClient();

  const { mutate } = useRenewContractMutation(client, {});

  const getStatusDisplay = useMemo(() => {
    let icon, text;
    switch (status) {
      case ContractStatus.Live:
        icon = <XSquare className='text-gray-500 mr-1' />;
        text = 'End contract...';
        break;
      case ContractStatus.Draft:
      case ContractStatus.Ended:
        icon = <DotLive className='text-gray-500 mr-1' />;
        text = 'Make live';
        break;
      case ContractStatus.OutOfContract:
        icon = <RefreshCw05 className='text-gray-500 mr-2' />;
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

  const handleChangeStatus = () => {
    switch (status) {
      case ContractStatus.Live:
        onOpenEndModal();
        break;
      case ContractStatus.Draft:
      case ContractStatus.Ended:
        onOpenStartModal();
        break;
      case ContractStatus.OutOfContract:
        mutate({ contractId });
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
        isEndModalOpen={isOpen}
        onHandleStatusChange={handleChangeStatus}
      />
      <ContractMenu
        onOpenEditModal={onOpenEditModal}
        statusContent={getStatusDisplay}
        onHandleStatusChange={handleChangeStatus}
        status={status}
      />
      <ContractEndModal
        isOpen={isOpen}
        onClose={onCloseEndModal}
        contractId={contractId}
        organizationName={organizationName}
        renewsAt={renewsAt}
        serviceStarted={serviceStarted}
        onUpdateContract={onUpdateContract}
        nextInvoiceDate={nextInvoiceDate}
      />
      <ContractStartModal
        isOpen={isStartModalOpen}
        onClose={onCloseStartModal}
        contractId={contractId}
        organizationName={organizationName}
        serviceStarted={serviceStarted}
        onUpdateContract={onUpdateContract}
      />
    </div>
  );
};
