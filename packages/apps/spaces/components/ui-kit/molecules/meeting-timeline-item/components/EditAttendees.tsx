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

export const PreviewAttendees = ({
  selectedContacts = [],
  onSelectContact = () => null,
  hiddenContactsNumber = 0,
}) => {
  const [dropdownOpen, setDropdownOpen] = useState<boolean>(false);

  const contactAutocompleteWrapperRef = useRef(null);

  useDetectClickOutside(contactAutocompleteWrapperRef, () => {
    console.log('🏷️ CLICKED OUTSIDE');
    setDropdownOpen(false);
  });

  return (
    <div ref={contactAutocompleteWrapperRef} style={{ position: 'relative' }}>
      <IconButton
        mode='secondary'
        onClick={() => setDropdownOpen(!dropdownOpen)}
        icon={<span style={{ fontSize: 12 }}>+{hiddenContactsNumber}</span>}
        size='xxxs'
        style={{
          width: 30,
          height: 30,
          zIndex: 3,
          // left: -60,
          position: 'relative',
        }}
      />

      {dropdownOpen && (
        <div
          className={classNames(
            styles.contactAutocompleteWrapper,
            styles.right,
          )}
        >
          {selectedContacts.map((contact) => {
            console.log('🏷️ ----- contact: ', contact);
            return (
              <li
                key={`contact-suggestion-${contact.id}`}
                className={classNames(styles.suggestionItem, styles.selected)}
                onClick={() => console.log()}
              >
                <ContactAvatar size={20} contactId={contact.id} showName />
                {/*<DeleteIconButton*/}
                {/*  onDelete={() => console.log('ON DELETE FROM ATTENDEES')}*/}
                {/*/>*/}
              </li>
            );
          })}
        </div>
      )}
    </div>
  );
};
