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
}

export const TimelineToolbelt: React.FC<ToolbeltProps> = ({
  onLogPhoneCall,
  onCreateNote,
  onCreateMeeting,
}) => {
  const { identity: loggedInUserEmail } = useRecoilValue(userData);

  return (
    <article className={styles.toolbelt}>
      <button
        aria-label='Create meeting'
        className={classNames(styles.button, styles.meeting)}
        onClick={onCreateMeeting}
      >
        <MeetingTimeline width={200} />
      </button>

      {onLogPhoneCall && (
        <button
          aria-label='Manually log phone call'
          className={styles.button}
          onClick={() =>
            onLogPhoneCall({
              appSource: 'Openline',
              sentBy: loggedInUserEmail,
              content: '',
              contentType: 'text/html',
            })
          }
        >
          <PhoneCallTimeline width={200} />
        </button>
      )}

      <button
        aria-label='Create Note'
        className={styles.button}
        onClick={() => onCreateNote({ appSource: 'Openline', html: '' })}
      >
        <NoteTimeline width={200} />
      </button>

      <div className={styles.belt} />
    </article>
  );
};
