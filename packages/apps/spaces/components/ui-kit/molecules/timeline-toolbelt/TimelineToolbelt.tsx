import React, { ReactNode } from 'react';
import { MeetingTimeline, NoteTimeline, PhoneCallTimeline } from '../../atoms';
import styles from './timeline-toolbelt.module.scss';
import classNames from 'classnames';
import {
  EditorMode,
  editorMode,
  showLegacyEditor,
} from '../../../../state/editor';
import { useSetRecoilState } from 'recoil';

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
  showPhoneCallButton?: boolean;
}

export const TimelineToolbelt: React.FC<ToolbeltProps> = ({
  onCreateNote,
  onCreateMeeting,
  isSkewed,
  showPhoneCallButton,
}) => {
  const setEditorMode = useSetRecoilState(editorMode);
  const setShowLegacyEditor = useSetRecoilState(showLegacyEditor);
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

      {showPhoneCallButton && (
        <button
          aria-label='Manually log phone call'
          className={classNames(styles.button, {
            [styles.isSkewed]: isSkewed,
          })}
          onClick={() => {
            setEditorMode({ mode: EditorMode.PhoneCall });
            setShowLegacyEditor(true);
          }}
        >
          <PhoneCallTimeline width={isSkewed ? 200 : 130} />
        </button>
      )}

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
