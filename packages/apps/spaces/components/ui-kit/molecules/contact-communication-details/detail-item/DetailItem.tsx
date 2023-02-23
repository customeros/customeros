import React, { useRef, useState } from 'react';
import Image from 'next/image';
import { IconButton } from '../../../atoms';
import classNames from 'classnames';
import styles from '../contact-communication-details.module.scss';
import { useDetectClickOutside } from '../../../../../hooks';
import { DetailItemEditMode } from './DetailItemEditMode';

export const DetailItem = ({
  id,
  isPrimary,
  label,
  data,
  // onChange,
  onDelete,
}: any) => {
  const listItemRef = useRef(null);
  const [editMode, setEditMode] = useState<boolean>(false);
  const [showButtons, setShowButtons] = useState<boolean>(false);
  useDetectClickOutside(listItemRef, () => {
    setShowButtons(false);
  });

  if (editMode) {
    return (
      <DetailItemEditMode
        id={id}
        isPrimary={isPrimary}
        label={label}
        data={data}
        onExitEditMode={() => {
          setShowButtons(false);
          setEditMode(false);
        }}
      />
    );
  }

  return (
    <li
      className={classNames(styles.communicationItem, {
        [styles.primary]: isPrimary,
        [styles.selected]: showButtons,
      })}
    >
      <div
        ref={listItemRef}
        role='button'
        tabIndex={0}
        key={id}
        className={styles.listContent}
        onClick={() => {
          setShowButtons(true);
        }}
      >
        <div className={styles.label}>
          <span>{label.toLowerCase()}</span>
        </div>
        <div className={styles.info}>
          <span> {data}</span>
        </div>
      </div>
      {showButtons && (
        <div className={styles.editButton}>
          <IconButton
            size='xxxs'
            mode='text'
            onClick={() => {
              setEditMode(true);
            }}
            icon={
              <Image
                src='/icons/pencil.svg'
                alt=''
                width={15}
                height={15}
                color='white'
              />
            }
          />
          <IconButton
            size='xxxs'
            mode='text'
            onClick={onDelete}
            icon={
              <Image
                src='/icons/trash.svg'
                alt=''
                width={15}
                height={15}
                color='white'
              />
            }
          />
        </div>
      )}
    </li>
  );
};
