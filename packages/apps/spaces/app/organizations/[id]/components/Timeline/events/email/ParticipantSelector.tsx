import React, { useRef, useState } from 'react';
import {
  useContactList,
  ComparisonOperator,
} from '@spaces/hooks/useContactList';
import { AutoComplete } from 'primereact/autocomplete';
import styles from './email-fields.module.scss';
import classNames from 'classnames';
import {Contact} from "@spaces/graphql";


export const EmailFields: React.FC = () => {
  const { onLoadContactList } = useContactList();
  const [filteredEmails, setFilteredEmails] = useState([]);
  const inputRef = useRef(null);

  const getContactSuggestions = async (filter: string) => {
    const response = await onLoadContactList({
      variables: {
        pagination: { page: 1, limit: 10 },
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
      const options = response?.data?.contacts?.content
        .filter((e: Contact) => e.emails.length)
        .map((e: Contact) => e.emails.map((email) => email.email))
        .flat();
      setFilteredEmails(options || []);
    }
  };
  // const handleKeyDown = (e: any) => {
  //   const {
  //     key,
  //     target: { value },
  //   } = e;
  //   switch (key) {
  //     case 'Tab':
  //       if (value) e.preventDefault();
  //       break;
  //     case 'Enter':
  //     case ',':
  //       {
  //         const trimmedValue = value.trim();
  //         if (trimmedValue) {
  //           setMode({
  //             ...editorModeState,
  //             to: [...editorModeState.to, trimmedValue],
  //           });
  //         }
  //         if (inputRef?.current) {
  //           // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  //           //@ts-expect-error
  //           inputRef.current.value = '';
  //         }
  //       }
  //       break;
  //     case 'Backspace':
  //       if (!value) {
  //         if (editorModeState.to.length > 0) {
  //           const newSelected = editorModeState.to.slice(0, -1);
  //
  //           setMode({
  //             ...editorModeState,
  //             to: newSelected,
  //           });
  //         }
  //       }
  //       break;
  //   }
  // };
  //
  // const handleSelect = (email: string) => {
  //   setMode({
  //     ...editorModeState,
  //     to: [...editorModeState.to, email],
  //   });
  // };
  // const handleChangeSubject = (subject: string) => {
  //   setMode({
  //     ...editorModeState,
  //     subject,
  //   });
  // };

  return (
    <div>
      <div className={`${styles.emailAutocomplete} emailEditor`}>
        <label className={styles.label}>To:</label>
        <div className={styles.autocompleteWithResults}>
          {/*<ul className={styles.selectedEmailsList}>*/}
          {/*  {editorModeState.to.map((e) => (*/}
          {/*    <li key={`selected-email-${e}`} className={styles.selectedEmail}>*/}
          {/*      {e}*/}
          {/*    </li>*/}
          {/*  ))}*/}
          {/*</ul>*/}
          <AutoComplete
            inputRef={inputRef}
            field='name'
            multiple
            value={''}
            suggestions={filteredEmails}
            itemTemplate={(email: string) => {
              return (
                <span
                  className={styles.option}
                  onClick={() => handleSelect(email)}
                >
                  {email}
                </span>
              );
            }}
            completeMethod={(e) => getContactSuggestions(e.query)}
            onChange={(e) => handleSelect(e.value)}
            onKeyDown={handleKeyDown}
          />
        </div>
      </div>
      <div className={styles.subjectWrapper} style={{ margin: 0 }}>
        <label className={classNames(styles.label, styles.subject)}>
          Subject:
          <input
            className={styles.subjectInput}
            onChange={({ target }) => handleChangeSubject(target.value)}
            value={editorModeState.subject}
            disabled
          />
        </label>
      </div>
    </div>
  );
};
