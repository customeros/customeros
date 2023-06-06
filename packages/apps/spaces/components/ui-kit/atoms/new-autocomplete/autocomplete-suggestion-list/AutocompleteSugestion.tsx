import React, { KeyboardEventHandler, MouseEventHandler } from 'react';
import styles from '../autocomplete.module.scss';
import classNames from 'classnames';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';

type Props = {
  item: SuggestionItem;
  active: boolean;
  onKeyDown?: KeyboardEventHandler<HTMLDivElement>;
  onClick: MouseEventHandler<HTMLLIElement>;
};

export const AutocompleteSuggestion = ({ item, active, onClick }: Props) => {
  return (
    <li
      tabIndex={0}
      className={classNames(styles.list_item, {
        [styles.active]: active,
      })}
      role='button'
      onClick={onClick}
    >
      <div className={styles.list_item_text}>{item.label}</div>
    </li>
  );
};
