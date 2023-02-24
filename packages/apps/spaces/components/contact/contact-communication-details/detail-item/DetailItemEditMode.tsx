import React, { ChangeEventHandler, useRef } from 'react';
import {
  Button,
  CheckSquare,
  DebouncedInput,
  Stop,
} from '../../../ui-kit/atoms';
import classNames from 'classnames';
import styles from '../contact-communication-details.module.scss';
import { useDetectClickOutside } from '../../../../hooks';
import { OverlayPanelEventType } from 'primereact';
import { EmailLabel } from '../../../../graphQL/generated';
import Image from 'next/image';
import { OverlayPanel } from '../../../ui-kit/atoms/overlay-panel';

interface Props {
  onChange?: ChangeEventHandler<HTMLInputElement>;
  onChangeLabelAndPrimary: (e: { label?: string; primary?: boolean }) => void;
  id: string;
  label: string;
  value: string;
  isPrimary: boolean;
  mode?: 'ADD' | 'EDIT';
  onExitEditMode: () => void;
}
export const DetailItemEditMode = ({
  id,
  isPrimary,
  mode = 'EDIT',
  label,
  value,
  onChange,
  onChangeLabelAndPrimary,
  onExitEditMode,
}: Props) => {
  const editListItemRef = useRef(null);
  const addCommunicationChannelContainerRef = useRef(null);

  useDetectClickOutside(editListItemRef, () => {
    onExitEditMode();
  });

  const labelOptions = Object.values(EmailLabel).map((labelOption) => ({
    label: labelOption.toLowerCase(),
    command: () => onChangeLabelAndPrimary({ label: labelOption }),
  }));

  return (
    <li
      key={id}
      ref={editListItemRef}
      className={classNames(styles.communicationItemEdit, {
        [styles.primary]: isPrimary,
      })}
    >
      <div className={styles.label}>
        <Button
          mode='link'
          style={{ display: 'inline-flex', paddingTop: 0, paddingBottom: 0 }}
          onClick={(e: OverlayPanelEventType) =>
            //@ts-expect-error revisit later
            addCommunicationChannelContainerRef?.current?.toggle(e)
          }
        >
          <div className={styles.editLabelIcon}>
            <Image
              src='/icons/code.svg'
              alt={'Change label'}
              height={12}
              width={12}
            />
            {label.toLowerCase()}
          </div>
        </Button>
        <OverlayPanel
          ref={addCommunicationChannelContainerRef}
          model={labelOptions}
        />
      </div>
      <div style={{ display: 'flex', flex: '1' }}>
        {mode === 'ADD' && onChange ? (
          <DebouncedInput
            onChange={onChange}
            inputSize='xxxs'
            value={value}
            placeholder='fill in your data'
            onKeyDown={(e) => e.key === 'Enter' && onExitEditMode()}
            debounceTimeout={0}
          />
        ) : (
          <div className={styles.info}>
            <span> {value}</span>
          </div>
        )}

        <Button
          mode='text'
          className={styles.primaryButton}
          style={{
            display: 'inline-flex',
            padding: 0,
            fontWeight: 'normal',
          }}
          onClick={() => onChangeLabelAndPrimary({ primary: !isPrimary })}
        >
          <div className={styles.editLabelIcon}>
            {isPrimary ? (
              <CheckSquare style={{ transform: 'scale(0.8)' }} />
            ) : (
              <Stop style={{ transform: 'scale(0.8)' }} />
            )}

            <span>Primary</span>
          </div>
        </Button>
      </div>
    </li>
  );
};
