import { useRef, useState, useEffect, KeyboardEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { EmailValidationDetails } from '@graphql/types';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface EmailCellProps {
  email: string;
  contactId: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailCell = observer(
  ({ email, validationDetails, contactId }: EmailCellProps) => {
    const emailInputRef = useRef<HTMLInputElement | null>(null);
    const store = useStore();
    const [isHovered, setIsHovered] = useState(false);

    const contactStore = store.contacts.value.get(contactId);

    const [isEdit, setIsEdit] = useState(false);
    const ref = useRef(null);

    useOutsideClick({
      ref: ref,
      handler: () => {
        setIsEdit(false);
      },
    });

    useEffect(() => {
      if (isHovered && isEdit) {
        emailInputRef.current?.focus();
      }
    }, [isHovered, isEdit]);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit]);

    const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
      if (e.key === 'Escape' || e.key === 'Enter') {
        emailInputRef?.current?.blur();
        setIsEdit(false);
      }
    };

    return (
      <div
        ref={ref}
        onKeyDown={handleEscape}
        className='flex justify-between'
        onDoubleClick={() => setIsEdit(true)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        {!isEdit && !email && <p className='text-gray-400'>Unknown</p>}
        {!isEdit && email && (
          <p className='overflow-ellipsis overflow-hidden'>{email}</p>
        )}
        {isEdit && (
          <Input
            size='xs'
            variant='unstyled'
            ref={emailInputRef}
            placeholder='Email'
            onFocus={(e) => e.target.select()}
            value={contactStore?.value?.emails?.[0]?.email ?? ''}
            onBlur={() => {
              if (!contactStore?.value?.emails?.[0]?.id) {
                contactStore?.addEmail();
              } else {
                contactStore?.updateEmail();
              }
            }}
            onChange={(e) => {
              contactStore?.update(
                (value) => {
                  set(value, 'emails[0].email', e.target.value);

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
        {email && (
          <EmailValidationMessage
            email={email}
            validationDetails={validationDetails}
          />
        )}
      </div>
    );
  },
);
