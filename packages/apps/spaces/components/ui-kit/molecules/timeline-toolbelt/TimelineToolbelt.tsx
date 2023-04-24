import React, { useState } from 'react';
import { MeetingTimeline, NoteTimeline, PhoneCallTimeline } from '../../atoms';
import styles from './timeline-toolbelt.module.scss';
import classNames from 'classnames';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../../state';

interface ToolbeltProps {
  onCreateMeeting: () => void;
  onCreateNote: ({
    appSource,
    html,
  }: {
    appSource: string;
    html: string;
  }) => void;
  isSkewed?: boolean;
  id?: string;
}

export const TimelineToolbelt: React.FC<ToolbeltProps> = ({
  onCreateNote,
  onCreateMeeting,
  isSkewed,
  id,
}) => {
  const { identity: loggedInUserEmail } = useRecoilValue(userData);
  const [deleteConfirmationModalVisible, setLogPhoneCallEditorVisible] =
    useState(false);
  return (
    <article
      className={classNames(styles.toolbelt, {
        [styles.isColumn]: !isSkewed,
      })}
    >
      <button
        aria-label='Create meeting'
        className={classNames(styles.button, {
          [styles.isSkewed]: isSkewed,
        })}
        onClick={onCreateMeeting}
      >
        <MeetingTimeline width={isSkewed ? 200 : 130} />
      </button>

      <button
        aria-label='Manually log phone call'
        className={classNames(styles.button, {
          [styles.isSkewed]: isSkewed,
        })}
        onClick={() => setLogPhoneCallEditorVisible(true)}
      >
        <PhoneCallTimeline width={200} />
      </button>

      <button
        aria-label='Create Note'
        className={classNames(styles.button, {
          [styles.isSkewed]: isSkewed,
        })}
        onClick={() => onCreateNote({ appSource: 'Openline', html: '' })}
      >
        <NoteTimeline width={isSkewed ? 200 : 130} />
      </button>

      {/*<Dialog*/}
      {/*  header={'Log phone call'}*/}
      {/*  draggable={false}*/}
      {/*  className={styles.dialog}*/}
      {/*  visible={deleteConfirmationModalVisible}*/}
      {/*  onHide={() => setLogPhoneCallEditorVisible(false)}*/}
      {/*>*/}
      {/*  <ContactEditor contactId={id} />*/}
      {/*</Dialog>*/}
    </article>
  );
};
