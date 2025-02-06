import { useRef } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
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
      <div ref={ref} className='flex group/contact-preview items-center'>
        {!contactName && (
          <p
            className='text-gray-700 font-medium no-underline hover:no-underline cursor-pointer'
            onClick={() => {
              if (store.ui.showPreviewCard && store.ui.focusRow === contactId) {
                store.ui.setShowPreviewCard(false);
              } else {
                store.ui.setFocusRow(contactId);
                store.ui.setShowPreviewCard(true);
              }
            }}
          >
            {isEnriching ? 'Enriching...' : 'Unnamed'}
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
              if (store.ui.showPreviewCard && store.ui.focusRow === contactId) {
                store.ui.setShowPreviewCard(false);
              } else {
                store.ui.setFocusRow(contactId);
                store.ui.setShowPreviewCard(true);
              }
            }}
          >
            {contactName}
          </p>
        )}
        <Icon
          name={'eye'}
          className='text-gray-400 ml-2 opacity-0 group-hover/contact-preview:opacity-100 cursor-pointer'
          onClick={() => {
            if (store.ui.showPreviewCard && store.ui.focusRow === contactId) {
              store.ui.setShowPreviewCard(false);
            } else {
              store.ui.setFocusRow(contactId);
              store.ui.setShowPreviewCard(true);
            }
          }}
        />
      </div>
    );
  },
);
