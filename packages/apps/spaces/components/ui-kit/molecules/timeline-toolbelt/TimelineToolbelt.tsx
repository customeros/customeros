import React from 'react';
import { MeetingTimeline, NoteTimeline, PhoneCallTimeline } from '../../atoms';
import styles from './timeline-toolbelt.module.scss';
import { className } from 'jsx-dom-cjs';
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
  onLogPhoneCall?: (input: any) => void;
  isSkewed?: boolean;
}

export const TimelineToolbelt: React.FC<ToolbeltProps> = ({
  onLogPhoneCall,
  onCreateNote,
  onCreateMeeting,
  isSkewed,
}) => {
  const { identity: loggedInUserEmail } = useRecoilValue(userData);

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
        aria-label='Create Note'
        className={classNames(styles.button, {
          [styles.isSkewed]: isSkewed,
        })}
        onClick={() => onCreateNote({ appSource: 'Openline', html: '' })}
      >
        <NoteTimeline width={isSkewed ? 200 : 130} />
      </button>

      <div className={styles.belt} />
    </article>
  );
};
