import { useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';
import { TableCellTooltip } from '@ui/presentation/Table';

interface ContractCellProps {
  contractId: string;
}

export const ContractCell = observer(({ contractId }: ContractCellProps) => {
  const navigate = useNavigate();

  const store = useStore();

  const contract = store.contracts.value.get(contractId);
  const id = store.organizations
    .toArray()
    .find((e) => e.contracts.find((c) => c.metadata.id === contractId))?.id;

  const linkRef = useRef<HTMLParagraphElement>(null);
  const [searchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const search = searchParams.get('search');
  const [lastSearchForPreset, setLastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });

  const handleNavigate = () => {
    if (!id) return;

    const href = getHref(id);

    if (!href) return;

    if (preset) {
      setLastSearchForPreset({
        ...lastSearchForPreset,
        [preset]: search ?? '',
      });
    }
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
          {contract?.value?.contractName ?? 'Unknown'}
        </p>
      </span>
    </TableCellTooltip>
  );
});

function getHref(id: string) {
  return `/organization/${id}?tab=account`;
}
