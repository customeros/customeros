import React, { FC, useRef, useState } from 'react';
import {
  ComparisonOperator,
  useContactMentionSuggestionsList,
} from '../../../../../hooks/useContactList';
import {
  Contact,
  MeetingParticipant,
  MeetingParticipantInput,
  UserParticipant,
} from '../../../../../graphQL/__generated__/generated';
import { getContactDisplayName } from '../../../../../utils';
import styles from './attendee-autocomplete.module.scss';
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

interface ContactAutocompleteProps {
  selectedAttendees: Array<MeetingParticipant>;
  onAddAttendee: (participantInput: MeetingParticipantInput) => void;
  onRemoveAttendee: (id: string) => void;
}

export const ContactAutocomplete: FC<ContactAutocompleteProps> = ({
  selectedAttendees = [],
  onAddAttendee,
  onRemoveAttendee,
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
    if (!dropdownOpen) return;
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
                  onClick={() =>
                    onAddAttendee({ contactID: value, type: 'contact' })
                  }
                  role='button'
                  tabIndex={0}
                >
                  <div>
                    <Avatar
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
            {!!selectedAttendees.length && (
              <li className={styles.listDivider}>Selected attendees:</li>
            )}

            {selectedAttendees.map((attendeeData) => {
              if (
                attendeeData.__typename !== 'ContactParticipant' &&
                attendeeData.__typename !== 'UserParticipant'
              )
                return null;

              const attendee =
                attendeeData.__typename === 'ContactParticipant'
                  ? attendeeData.contactParticipant
                  : (attendeeData as UserParticipant).userParticipant;

              return (
                <li
                  key={`contact-suggestion-${attendee.id}`}
                  className={classNames(styles.suggestionItem, styles.selected)}
                >
                  <ContactAvatar size={20} contactId={attendee.id} showName />
                  <DeleteIconButton
                    onDelete={() => onRemoveAttendee(attendee.id)}
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
