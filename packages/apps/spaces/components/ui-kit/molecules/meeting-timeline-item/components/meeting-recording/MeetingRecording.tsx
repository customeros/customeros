import React from 'react';
import styles from './meeting-timeline-item.module.scss';

import classNames from 'classnames';
import {
  IconButton,
  CloudUpload,
  VoiceWaveRecord,
  ChevronDown,
} from '../../../../atoms';
import FileO from '../../../../atoms/icons/FileO';
import { Meeting } from '../../../../../../hooks/useMeeting';
import { useFileUpload } from '../../../../../../hooks/useFileUpload';
import { toast } from 'react-toastify';

interface MeetingTimelineItemProps {
  meeting: Meeting;
  onUpdateMeetingRecording: (id: string | null) => void;
}

export const MeetingRecording = ({
  meeting,
  onUpdateMeetingRecording,
}: MeetingTimelineItemProps): JSX.Element => {
  const uploadInputRef = React.useRef<HTMLInputElement>(null);

  const { isDraggingOver, handleDrag, handleDrop, handleInputFileChange } =
    useFileUpload({
      prevFiles: [],
      onBeginFileUpload: (data) => console.log('onBeginFileUpload', data),
      onFileUpload: (data) => onUpdateMeetingRecording(data.id),
      onFileUploadError: () =>
        toast.error(
          'Something went wrong while uploading recording of a meeting',
        ),
      onFileRemove: () => onUpdateMeetingRecording(null),
      uploadInputRef,
    });

  return (
    <>
      <article
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        className={classNames(styles.recordingSection, {
          [styles.recordingUploaded]: meeting.recording,
          [styles.isDraggingOver]: isDraggingOver,
        })}
      >
        <div
          className={classNames(styles.recordingCta, {
            [styles.recordingUploaded]: meeting.recording,
          })}
          onClick={() => uploadInputRef?.current?.click()}
        >
          <input
            style={{ display: 'none' }}
            ref={uploadInputRef}
            type='file'
            onChange={handleInputFileChange}
          />

          {meeting.recording ? (
            <div className={styles.recordingIcon}>
              <FileO height={24} width={24} aria-label='Meeting recording' />
            </div>
          ) : (
            <div className={styles.recordingIcon} style={{ padding: 2 }}>
              <CloudUpload height={20} width={20} />
            </div>
          )}
          {meeting.recording ? (
            <span> Meeting recording uploaded </span>
          ) : (
            <h3>
              {isDraggingOver ? 'Drop file here' : 'Upload the recording'}{' '}
            </h3>
          )}
          <VoiceWaveRecord />
        </div>
      </article>

      <div className={styles.collapsibleSection}>
        <div
          className={classNames(styles.transcriptionSection, {
            [styles.recordingUploaded]: meeting.recording,
            [styles.isDraggingOver]: isDraggingOver,
          })}
        />
        <div className={styles.collapseExpandButtonWrapper}>
          <IconButton
            className={styles.collapseExpandButton}
            isSquare
            disabled
            mode='secondary'
            size='xxxxs'
            icon={<ChevronDown width={24} height={24} />}
            onClick={() => console.log('collapse / expand button click')}
          />
        </div>
      </div>

      <section>{/* collapsible section*/}</section>
    </>
  );
};
