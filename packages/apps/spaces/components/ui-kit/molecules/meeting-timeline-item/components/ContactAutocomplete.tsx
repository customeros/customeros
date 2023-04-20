import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  ComparisonOperator,
  useContactMentionSuggestionsList,
} from '../../../../../hooks/useContactList';
import {
  Contact,
  Organization,
} from '../../../../../graphQL/__generated__/generated';
import { MentionAtomState, MentionAtomNodeAttributes } from '@remirror/react';
import { getContactDisplayName } from '../../../../../utils';
import styles from './contact-autocomplere.module.scss';
import {
  Avatar,
  DebouncedInput,
  DeleteIconButton,
  Highlight,
  IconButton,
  Plus,
  Times,
  User,
  UserEdit,
} from '../../../atoms';
import { useCreateContact } from '../../../../../hooks/useContact';
import { toast } from 'react-toastify';
import classNames from 'classnames';
import { useDetectClickOutside } from '../../../../../hooks';
import { ContactAvatar } from '../../contact-avatar/ContactAvatar';

interface MentionProps<
  UserData extends MentionAtomNodeAttributes = MentionAtomNodeAttributes,
> {
  users?: UserData[];
  tags?: string[];
}

export const ContactAutocomplete = ({
  selectedContacts = [],
  onSelectContact = () => null,
}) => {
  const { onLoadContactMentionSuggestionsList } =
    useContactMentionSuggestionsList();
  const { onCreateContact } = useCreateContact();
  const [inputValue, setInputValue] = useState<string>('');
  const [dropdownOpen, setDropdownOpen] = useState<boolean>(false);
  const [filteredContacts, setFilteredContacts] = useState<
    Array<{ value: string; label: string }>
  >([]);

  const contactAutocompleteWrapperRef = useRef(null);

  useDetectClickOutside(contactAutocompleteWrapperRef, () => {
    console.log('🏷️ CLICKED OUTSIDE');
    setInputValue('');
    setDropdownOpen(false);
  });

  const getContactSuggestions = async (filter: string) => {
    if (!filter.length) {
      setFilteredContacts([]);
      return;
    }

    const response = await onLoadContactMentionSuggestionsList({
      variables: {
        pagination: { page: 0, limit: 10 },
        where: {
          OR: [
            {
              filter: {
                property: 'FIRST_NAME',
                value: filter,
                operation: ComparisonOperator.Contains,
              },
            },
            {
              filter: {
                property: 'LAST_NAME',
                value: filter,
                operation: ComparisonOperator.Contains,
              },
            },
          ],
        },
      },
    });
    if (response?.data) {
      const options = response?.data?.contacts?.content.map((e: Contact) => ({
        label: getContactDisplayName(e),
        value: e.id,
      }));
      setFilteredContacts(options || []);
    }
  };

  return (
    <div ref={contactAutocompleteWrapperRef}>
      <IconButton
        mode='secondary'
        onClick={() => setDropdownOpen(!dropdownOpen)}
        icon={<UserEdit />}
        size='xxxs'
      />

      {dropdownOpen && (
        <div
          className={classNames(
            styles.contactAutocompleteWrapper,
            styles.right,
          )}
        >
          <DebouncedInput
            inlineMode
            className={styles.contactAutocompleteInput}
            placeholder='Add attendees'
            onChange={(event) => {
              setInputValue(event.target.value);
              getContactSuggestions(event.target.value);
            }}
          />
          <ul>
            {filteredContacts.map(({ label, value }) => {
              const name = label.split(' ');

              return (
                <li
                  key={`contact-suggestion-${value}`}
                  className={classNames(
                    styles.suggestionItem,
                    styles.selectable,
                  )}
                  onClick={() => console.log(label)}
                  role='button'
                  tabIndex={0}
                >
                  <div>
                    <Avatar
                      onlyName
                      name={name?.[0] || ''}
                      surname={name.length === 2 ? name[1] : name[2]}
                      size={20}
                      image={
                        (name.length === 1 || name[0] === 'Unnamed') && (
                          <User style={{ transform: 'scale(0.6)' }} />
                        )
                      }
                    />
                  </div>

                  <span>
                    <Highlight text={label || ''} highlight={inputValue} />{' '}
                  </span>
                  <Plus style={{ transform: 'scale(0.6)' }} />
                </li>
              );
            })}
            {!!inputValue.length && !filteredContacts.length && (
              <li
                role='button'
                tabIndex={0}
                onClick={(e) => {
                  console.log('🏷️ ----- e: ADD NEW ', e);
                  const name = inputValue.split(' ');
                  if (name.length === 0) {
                    toast.error('Could not create contact with empty name');
                    return;
                  }
                  if (name.length === 1) {
                    return onCreateContact({ firstName: inputValue });
                  }
                  return onCreateContact({
                    firstName: name[0],
                    lastName: name[1],
                  });
                }}
              >
                Create contact &apos;{inputValue}&apos;
              </li>
            )}
            <li className={styles.listDivider}>Selected attendees:</li>
            {selectedContacts.map((contact) => {
              console.log('🏷️ ----- contact: ', contact);
              return (
                <li
                  key={`contact-suggestion-${contact.id}`}
                  className={classNames(styles.suggestionItem, styles.selected)}
                  onClick={() => console.log()}
                >
                  <ContactAvatar size={20} contactId={contact.id} onlyName />
                  <DeleteIconButton
                    onDelete={() => console.log('ON DELETE FROM ATTENDEES')}
                  />
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
};
