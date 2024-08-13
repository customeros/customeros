import { useRef } from 'react';
import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TableCellTooltip } from '@ui/presentation/Table';

export const ContractCell = observer(
  ({
    organizationId,
    contractId,
  }: {
    contractId: string;
    organizationId: string;
  }) => {
    const store = useStore();
    const itemRef = useRef<HTMLAnchorElement>(null);

    const organization = store.organizations?.value?.get(organizationId)?.value;
    const contract = store.contracts?.value?.get(contractId)?.value;
    const name = contract?.contractName || `${organization?.name}'s contract`;

    return (
      <TableCellTooltip
        hasArrow
        label={name}
        align='start'
        side='bottom'
        targetRef={itemRef}
      >
        <div className='overflow-hidden overflow-ellipsis'>
          <p className='text-xs text-gray-500'>{organization?.name}</p>
          <Link
            ref={itemRef}
            to={`/organization/${organizationId}?tab=account`}
            className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors overflow-hidden overflow-ellipsis block'
          >
            {name}
          </Link>
        </div>
      </TableCellTooltip>
    );
  },
);
