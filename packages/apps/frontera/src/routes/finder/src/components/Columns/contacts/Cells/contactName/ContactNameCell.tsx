import { useRef, useState, useEffect, KeyboardEvent } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';

interface ContactNameCellProps {
  contactId: string;
  canNavigate?: boolean;
}

export const ContactNameCell = observer(
  ({ contactId, canNavigate }: ContactNameCellProps) => {
    const contactNameInputRef = useRef<HTMLInputElement | null>(null);
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);

    const contactStore = store.contacts.value.get(contactId);
    const contactName = contactStore?.name;

    const [isEdit, setIsEdit] = useState(false);
    const ref = useRef(null);

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

    const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
      if (e.key === 'Escape' || e.key === 'Enter') {
        contactNameInputRef?.current?.blur();
        setIsEdit(false);
      }
      e.stopPropagation();
    };

    return (
      <div
        ref={ref}
        className='flex'
        onDoubleClick={() => setIsEdit(true)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        {!isEdit && !contactName && <p className='text-gray-400'>Unknown</p>}
        {!isEdit && contactName && (
          <p
            role='button'
            data-test={`contact-name-in-contacts-table`}
            className={cn(
              'overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer',
              !canNavigate && 'cursor-default',
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
        {isEdit && (
          <Input
            size='xs'
            placeholder='Name'
            variant='unstyled'
            onKeyDown={handleEscape}
            ref={contactNameInputRef}
            value={contactStore?.name ?? ''}
            onFocus={(e) => e.target.select()}
            className={'font-medium placeholder-font-normal'}
            onBlur={() => {
              contactStore?.updateContactName();
            }}
            onChange={(e) => {
              contactStore?.update(
                (value) => {
                  value.name = e.target.value;

                  return value;
                },
                { mutate: false },
              );
            }}
          />
        )}

        {isHovered && !isEdit && (
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='edit'
            className='ml-3 rounded-[5px]'
            onClick={() => setIsEdit(!isEdit)}
            icon={<Edit03 className='text-gray-500' />}
          />
        )}
      </div>
    );
  },
);
