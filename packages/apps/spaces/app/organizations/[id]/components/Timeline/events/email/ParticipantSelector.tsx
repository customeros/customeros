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
  const [filteredEmails, setFilteredEmails] = useState([]);
  const inputRef = useRef(null);


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
