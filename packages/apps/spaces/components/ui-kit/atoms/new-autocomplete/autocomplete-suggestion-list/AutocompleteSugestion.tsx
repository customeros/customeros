import React, { KeyboardEventHandler, MouseEventHandler, useRef } from 'react';
import styles from '../autocomplete.module.scss';
import classNames from 'classnames';
import { SuggestionItem } from '@spaces/atoms/new-autocomplete/Autocomplete';

type Props = {
  item: SuggestionItem;
  active: boolean;
  onKeyDown?: KeyboardEventHandler<HTMLDivElement>;
  onClick: MouseEventHandler<HTMLDivElement>;
};

export const AutocompleteSuggestion = ({ item, active, onClick }: Props) => {
  return (
    <div
      tabIndex={0}
      className={classNames(styles.list_item, {
        [styles.active]: active,
      })}
      role='listitem'
      onClick={onClick}
    >
      <div className={styles.list_item_text}>{item.label}</div>
    </div>
  );
};
