import { useNavigate } from 'react-router-dom';
import React, { useRef, useState, useEffect, KeyboardEvent } from 'react';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';

interface ContactNameCellProps {
  contactId: string;
}

export const ContactNameCell: React.FC<ContactNameCellProps> = observer(
  ({ contactId }) => {
    const contactNameInputRef = useRef<HTMLInputElement | null>(null);
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);

    const contactStore = store.contacts.value.get(contactId);
    const contactName = contactStore?.value.name;
    const navigate = useNavigate();
    const [tabs] = useLocalStorage<{
      [key: string]: string;
    }>(`customeros-player-last-position`, { root: 'organization' });

    const lastPositionParams = contactStore?.organizationId
      ? tabs[contactStore?.organizationId]
      : undefined;
    const href = contactStore?.organizationId
      ? getHref(
          contactStore?.organizationId,
          'tab=people' || lastPositionParams,
        )
      : null;

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
    };

    return (
      <div
        onDoubleClick={() => setIsEdit(true)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        ref={ref}
        className='flex'
        onKeyDown={handleEscape}
      >
        {!isEdit && !contactName && <p className='text-gray-400'>Unknown</p>}
        {!isEdit && contactName && (
          <p
            role='button'
            className='max-w-[140px] overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer'
            onClick={() => href && navigate(href)}
          >
            {contactName}
          </p>
        )}
        {isEdit && (
          <Input
            ref={contactNameInputRef}
            onFocus={(e) => e.target.select()}
            placeholder='Name'
            variant='unstyled'
            className={'font-medium placeholder-font-normal'}
            size='xs'
            value={contactStore?.value?.name ?? ''}
            onChange={(e) => {
              contactStore?.update(
                (value) => {
                  value.name = e.target.value;

                  return value;
                },
                { mutate: false },
              );
            }}
            onBlur={(e) => {
              contactStore?.update((value) => {
                value.name = e.target.value;

                return value;
              });
            }}
          />
        )}

        {isHovered && !isEdit && (
          <IconButton
            className='ml-3 rounded-[5px]'
            variant='ghost'
            size='xxs'
            onClick={() => setIsEdit(!isEdit)}
            aria-label='edit'
            icon={<Edit03 className='text-gray-500' />}
          />
        )}
      </div>
    );
  },
);
function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
