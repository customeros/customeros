import React, { useState } from 'react';
import { MeetingTimeline, NoteTimeline } from '../../atoms';
import styles from './timeline-toolbelt.module.scss';
import classNames from 'classnames';

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
}

export const TimelineToolbelt: React.FC<ToolbeltProps> = ({
  onCreateNote,
  onCreateMeeting,
  isSkewed,
}) => {
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
    </article>
  );
};
