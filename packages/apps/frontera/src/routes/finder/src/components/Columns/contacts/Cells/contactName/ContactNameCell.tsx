import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';

interface ContactNameCellProps {
  contactId: string;
  canNavigate?: boolean;
}

export const ContactNameCell = observer(
  ({ contactId, canNavigate }: ContactNameCellProps) => {
    const contactNameInputRef = useRef<HTMLInputElement | null>(null);
    const store = useStore();

    const contactStore = store.contacts.value.get(contactId);
    const contactName = contactStore?.name;

    const [isEdit, setIsEdit] = useState(false);
    const ref = useRef(null);

    const isEnriching = contactStore?.isEnriching;

    useOutsideClick({
      ref: ref,
      handler: () => {
        setIsEdit(false);
      },
    });

    useEffect(() => {
      if (isEdit) {
        contactNameInputRef.current?.focus();
      }
    }, [isEdit]);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit]);

    if (!contactStore) return;

    return (
      <div ref={ref} className='flex'>
        {!isEdit && !contactName && (
          <p className='text-gray-400'>
            {isEnriching ? 'Enriching...' : 'Unknown'}
          </p>
        )}
        {!isEdit && contactName && (
          <p
            role='button'
            data-test={`contact-name-in-contacts-table`}
            className={cn(
              'overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer',
              canNavigate && 'cursor-default',
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
