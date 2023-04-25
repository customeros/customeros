import React, { FC, useEffect, useRef, useState } from 'react';
import {
  ComparisonOperator,
  useContactMentionSuggestionsList,
} from '../../../../../hooks/useContactList';
import {
  Contact,
  MeetingParticipant,
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
  User,
  UserEdit,
} from '../../../atoms';
import { useCreateContact } from '../../../../../hooks/useContact';
import { toast } from 'react-toastify';
import classNames from 'classnames';
import { useDetectClickOutside } from '../../../../../hooks';
import { ContactAvatar } from '../../contact-avatar/ContactAvatar';
import { useUsers } from '../../../../../hooks/useUser';
import { useAutoAnimate } from '@formkit/auto-animate/react';
import {
  useLinkMeetingAttendee,
  useUnlinkMeetingAttendee,
} from '../../../../../hooks/useMeeting';

interface AttendeeAutocompleteProps {
  selectedAttendees: Array<MeetingParticipant>;
  meetingId: string;
}

export const AttendeeAutocomplete: FC<AttendeeAutocompleteProps> = ({
  selectedAttendees = [],
  meetingId,
}) => {
  const { onLoadContactMentionSuggestionsList } =
    useContactMentionSuggestionsList();
  const { onLoadUsers } = useUsers();
  const { onCreateContact } = useCreateContact();
  const [inputValue, setInputValue] = useState<string>('');
  const [dropdownOpen, setDropdownOpen] = useState<boolean>(false);
  const { onLinkMeetingAttendee } = useLinkMeetingAttendee({
    meetingId: meetingId,
  });
  const { onUnlinkMeetingAttendee } = useUnlinkMeetingAttendee({
    meetingId: meetingId,
  });
  const [filteredContacts, setFilteredContacts] = useState<
    Array<{ value: string; label: string; type?: string }>
  >([]);
  const [animateRef] = useAutoAnimate({
    easing: 'linear',
  });

  const attendeeAutocompleteWrapperRef = useRef(null);

  useDetectClickOutside(attendeeAutocompleteWrapperRef, () => {
    if (!dropdownOpen) return;
    setInputValue('');
    setDropdownOpen(false);
  });

  useEffect(() => {
    getContactSuggestions(inputValue);
  }, [inputValue]);

  const getContactSuggestions = async (filter: string) => {
    if (!filter.length) {
      setFilteredContacts([]);
      return;
    }

    const contactResponse = await onLoadContactMentionSuggestionsList({
      variables: {
        pagination: { page: 0, limit: 5 },
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

    const userResponse = await onLoadUsers({
      variables: {
        pagination: { page: 0, limit: 5 },
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
    if (contactResponse?.data || userResponse.data) {
      const options = [
        ...(contactResponse?.data?.contacts?.content || []),
        ...(userResponse?.data?.users?.content || []),
      ].map((e: Contact) => ({
        label: getContactDisplayName(e),
        value: e.id,
        type: e.__typename,
      }));
      setFilteredContacts(options || []);
    }
  };
  console.log('🏷️ ----- selectedAttendees: ', selectedAttendees);
  console.log('🏷️ ----- filteredContacts: ', filteredContacts);
  return (
    <div ref={attendeeAutocompleteWrapperRef}>
      <IconButton
        mode='secondary'
        onClick={() => setDropdownOpen(!dropdownOpen)}
        icon={<UserEdit />}
        size='xxxs'
      />

      {dropdownOpen && (
        <div
          className={classNames(
            styles.attendeeAutocompleteWrapper,
            styles.right,
          )}
        >
          <DebouncedInput
            inlineMode
            className={styles.attendeeAutocompleteInput}
            placeholder='Add attendees'
            onChange={(event) => {
              setInputValue(event.target.value);
            }}
          />
          <ul ref={animateRef}>
            {filteredContacts.map(({ label, value, type }) => {
              const name = label.split(' ');

              return (
                <li
                  key={`contact-suggestion-${value}`}
                  className={classNames(
                    styles.suggestionItem,
                    styles.selectable,
                  )}
                  onClick={() => {
                    const payload =
                      type === 'Contact'
                        ? { contactId: value }
                        : { userId: value };
                    onLinkMeetingAttendee(payload);

                    setInputValue('');
                  }}
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
                className={styles.listDivider}
                tabIndex={0}
                onClick={(e) => {
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
                    onDelete={() => {
                      const payload =
                        attendee.__typename === 'Contact'
                          ? { contactId: attendee.id }
                          : { userId: attendee.id };
                      onUnlinkMeetingAttendee(payload);
                    }}
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
