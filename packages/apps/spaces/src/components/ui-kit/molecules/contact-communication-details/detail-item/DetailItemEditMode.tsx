import React, { useRef } from 'react';
import { DebouncedInput } from '../../../atoms';
import classNames from 'classnames';
import styles from './../contact-communication-details.module.scss';
import { useDetectClickOutside } from '../../../../../hooks';

export const DetailItemEditMode = ({
  id,
  isPrimary,
  label,
  data,
  onChange,
  onExitEditMode,
}: any) => {
  const editListItemRef = useRef(null);
  useDetectClickOutside(editListItemRef, () => {
    onExitEditMode();
  });

  return (
    <li
      key={id}
      ref={editListItemRef}
      className={classNames(styles.communicationItemEdit, {
        [styles.primary]: isPrimary,
      })}
    >
      <div className={styles.label}>
        <span>{label}</span>
      </div>

      <DebouncedInput onChange={onChange} inputSize='xxxs' value={data} />
    </li>
  );
};
