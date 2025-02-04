import { useRef } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store';

import { cn } from '@ui/utils/cn.ts';
import { useStore } from '@shared/hooks/useStore';

interface ContractCellProps {
  contractId: string;
}

export const ContractCell = observer(({ contractId }: ContractCellProps) => {
  const navigate = useNavigate();

  const store = useStore();

  const contract = store.contracts.value.get(contractId) as ContractStore;
  const linkRef = useRef<HTMLParagraphElement>(null);

  const handleNavigate = () => {
    if (!contract?.organizationId) return;

    const href = getHref(contract.organizationId);

    if (!href) return;

    navigate(href);
  };

  return (
    <div
      role='button'
      ref={linkRef}
      onClick={handleNavigate}
      data-test='Contract-name-in-all-orgs-table'
      className={cn(
        'overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer pr-7',
      )}
    >
      {contract?.value?.contractName || `${contract.organization?.value?.name}`}
    </div>
  );
});

function getHref(id: string) {
  return `/organization/${id}?tab=account`;
}
