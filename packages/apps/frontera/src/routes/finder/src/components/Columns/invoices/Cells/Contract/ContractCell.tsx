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
    const orgName = organization?.name;

    return (
      <div className='overflow-hidden overflow-ellipsis'>
        {orgName ? (
          <Link
            ref={itemRef}
            to={`/organization/${organizationId}?tab=account`}
            className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors'
          >
            {name}
          </Link>
        ) : (
          <TableCellTooltip
            hasArrow
            align='start'
            side='bottom'
            targetRef={itemRef}
            label='The org linked to this contract does not exist'
          >
            <div>
              <span className='text-gray-700 font-medium cursor-pointer'>
                {name}
              </span>
            </div>
          </TableCellTooltip>
        )}
      </div>
    );
  },
);
