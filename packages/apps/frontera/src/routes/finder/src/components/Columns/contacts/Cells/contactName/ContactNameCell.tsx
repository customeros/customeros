import { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';

interface ContactNameCellProps {
  contactId: string;
}

export const ContactNameCell = observer(
  ({ contactId }: ContactNameCellProps) => {
    const store = useStore();

    const contactStore = store.contacts.value.get(contactId);
    const contactName = contactStore?.name;

    const ref = useRef(null);

    const isEnriching = contactStore?.isEnriching;

    if (!contactStore) return;

    return (
      <div ref={ref} className='flex'>
        {!contactName && (
          <p className='text-gray-400'>
            {isEnriching ? 'Enriching...' : 'Unknown'}
          </p>
        )}
        {contactName && (
          <p
            role='button'
            data-test={`contact-name-in-contacts-table`}
            className={cn(
              'overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline',
            )}
            onClick={() => {
              if (
                store.ui.contactPreviewCardOpen === true &&
                store.ui.focusRow === contactId
              ) {
                store.ui.setContactPreviewCardOpen(false);
              } else {
                store.ui.setFocusRow(contactId);
                store.ui.setContactPreviewCardOpen(true);
              }
            }}
          >
            {contactName}
          </p>
        )}
      </div>
    );
  },
);
