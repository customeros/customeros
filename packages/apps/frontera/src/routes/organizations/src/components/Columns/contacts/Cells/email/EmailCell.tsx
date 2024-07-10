import React, { useRef, useState, useEffect, KeyboardEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { EmailValidationDetails } from '@graphql/types';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { SimpleValidationIndicator } from '@ui/presentation/validation/simple-validation-indicator';
import { VALIDATION_MESSAGES } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/utils.ts';
function isValidEmail(email: string) {
  // Regular expression for validating an email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  // Test the email against the regex
  return emailRegex.test(email);
}

interface EmailCellProps {
  email: string;
  contactId: string;
  validationDetails: EmailValidationDetails | undefined;
}

export const EmailCell: React.FC<EmailCellProps> = observer(
  ({ email, validationDetails, contactId }) => {
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

    const getMessages = () => {
      if (!validationDetails) return [];
      const { validated, isReachable, isValidSyntax } = validationDetails;
      if (validated && !isValidEmail(email) && isReachable === 'safe')
        return [VALIDATION_MESSAGES.isValidSyntax.message];

      if (!validated && !isValidSyntax && isReachable !== 'safe') {
        return [VALIDATION_MESSAGES.isValidSyntax.message];
      }
      if (
        validated &&
        isReachable &&
        (VALIDATION_MESSAGES.isReachable.condition as Array<string>).includes(
          isReachable,
        )
      ) {
        return [VALIDATION_MESSAGES.isReachable.message];
      }

      return [];
    };

    const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
      if (e.key === 'Escape' || e.key === 'Enter') {
        emailInputRef?.current?.blur();
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
        {!isEdit && !email && <p className='text-gray-400'>Unknown</p>}
        {!isEdit && email && (
          <p className='max-w-[140px] overflow-ellipsis overflow-hidden'>
            {email}
          </p>
        )}
        {isEdit && (
          <Input
            ref={emailInputRef}
            placeholder='Email'
            onFocus={(e) => e.target.select()}
            variant='unstyled'
            size='xs'
            value={contactStore?.value?.emails?.[0]?.email ?? ''}
            onChange={(e) => {
              contactStore?.update(
                (value) => {
                  set(value, 'emails[0].email', e.target.value);

                  return value;
                },
                { mutate: false },
              );
            }}
            onBlur={() => {
              if (!contactStore?.value?.emails?.[0]?.id) {
                contactStore?.addEmail();
              } else {
                contactStore?.updateEmail();
              }
            }}
          />
        )}

        {email && (
          <SimpleValidationIndicator
            errorMessages={getMessages()}
            showValidationMessage={true}
            isLoading={false}
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
