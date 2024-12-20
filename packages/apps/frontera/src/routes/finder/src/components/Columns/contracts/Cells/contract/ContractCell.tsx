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
  const org = contract.organization;
  const linkRef = useRef<HTMLParagraphElement>(null);

  const handleNavigate = () => {
    if (!org?.id) return;

    const href = getHref(org?.id);

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
        {
          'text-gray-400 cursor-not-allowed': !org?.id,
        },
      )}
    >
      {contract?.value?.contractName || `${org?.value?.name}`}
    </div>
  );
});

function getHref(id: string) {
  return `/organization/${id}?tab=account`;
}
