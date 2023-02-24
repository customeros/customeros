import React, { ChangeEventHandler, useRef, useState } from 'react';
import Image from 'next/image';
import { IconButton } from '../../../ui-kit/atoms';
import classNames from 'classnames';
import styles from '../contact-communication-details.module.scss';
import { useDetectClickOutside } from '../../../../hooks';
import { DetailItemEditMode } from './DetailItemEditMode';
import { EmailLabel, PhoneNumberLabel } from '../../../../graphQL/generated';

interface Props {
  onChange?: ChangeEventHandler<HTMLInputElement>;
  onChangeLabelAndPrimary: (
    e: { label?: string } | { primary?: boolean },
  ) => void;
  id: string;
  label: string;
  data: string;
  isPrimary: boolean;
  mode?: 'ADD' | 'EDIT';
  onDelete: () => void;
  labelOptionEnum: typeof EmailLabel | typeof PhoneNumberLabel;
}
export const DetailItem = ({
  id,
  isPrimary,
  label,
  data,
  // onChange,
  onChangeLabelAndPrimary,
  onDelete,
  labelOptionEnum,
}: Props) => {
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
        value={data}
        onChangeLabelAndPrimary={onChangeLabelAndPrimary}
        labelOptionEnum={labelOptionEnum}
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
        onKeyDown={(e) => e.key === 'Enter' && setShowButtons(true)}
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
