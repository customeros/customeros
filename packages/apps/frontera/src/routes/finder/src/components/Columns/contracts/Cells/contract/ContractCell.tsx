import { useRef } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { TableCellTooltip } from '@ui/presentation/Table';

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
    <TableCellTooltip
      hasArrow
      align='start'
      side='bottom'
      targetRef={linkRef}
      label={contract?.value?.contractName ?? ''}
    >
      <span className='inline'>
        <p
          role='button'
          ref={linkRef}
          onClick={handleNavigate}
          data-test='Contract-name-in-all-orgs-table'
          className='overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer pr-7'
        >
          {contract?.value?.contractName || `${org?.value?.name}`}
        </p>
      </span>
    </TableCellTooltip>
  );
});

function getHref(id: string) {
  return `/organization/${id}?tab=account`;
}
