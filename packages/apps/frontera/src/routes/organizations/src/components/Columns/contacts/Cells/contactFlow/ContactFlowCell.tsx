import { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TableCellTooltip } from '@ui/presentation/Table';

interface ContactNameCellProps {
  contactId: string;
}

export const ContactFlowCell = observer(
  ({ contactId }: ContactNameCellProps) => {
    const store = useStore();

    const contactStore = store.contacts.value.get(contactId);
    const flowName = contactStore?.flow?.value?.name;
    const itemRef = useRef<HTMLDivElement>(null);

    if (!flowName) return <div className='text-gray-400'>None</div>;

    return (
      <TableCellTooltip
        hasArrow
        align='start'
        side='bottom'
        label={flowName}
        targetRef={itemRef}
      >
        <div ref={itemRef} className='flex overflow-hidden'>
          <div className=' overflow-x-hidden overflow-ellipsis'>{flowName}</div>
        </div>
      </TableCellTooltip>
    );
  },
);
