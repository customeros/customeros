import React, { useRef } from 'react';
import { Button, DebouncedInput } from '../../../atoms';
import classNames from 'classnames';
import styles from '../contact-communication-details.module.scss';
import { useDetectClickOutside } from '../../../../../hooks';
import { Dropdown, OverlayPanelEventType } from 'primereact';
import { EmailLabel } from '../../../../../graphQL/types';
import Image from 'next/image';
import { OverlayPanel } from '../../../atoms/overlay-panel';
import { log } from 'util';

// add label dropdown
// add label options
export const DetailItemEditMode = ({
  id,
  isPrimary,
  label,
  data,
  onChange,
  onExitEditMode,
}: any) => {
  const editListItemRef = useRef(null);
  const addCommunicationChannelContainerRef = useRef(null);

  useDetectClickOutside(editListItemRef, () => {
    onExitEditMode();
  });

  const labelOptions = Object.values(EmailLabel).map((label) => ({
    label: label.toLowerCase(),
    command: () => console.log('HEY'),
  }));

  console.log('🏷️ ----- la: ', labelOptions);

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
          role={'button'}
          tabIndex={0}
          mode='link'
          style={{ display: 'inline-flex', paddingTop: 0 }}
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

      <DebouncedInput
        onChange={onChange}
        inputSize='xxxs'
        value={data}
        placeholder='fill in your data'
      />
    </li>
  );
};
