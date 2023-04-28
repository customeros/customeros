import React, { useState } from 'react';
import classNames from 'classnames';
import { toast } from 'react-toastify';
import {
  IconButton,
  CloudUpload,
  VoiceWaveRecord,
  ChevronDown,
  Checkbox,
} from '../../../../atoms';
import FileO from '../../../../atoms/icons/FileO';
import { Meeting } from '../../../../../../hooks/useMeeting';
import { useFileUpload } from '../../../../../../hooks/useFileUpload';
import { Message } from '../../../../atoms/message/Message';
import styles from './meeting-timeline-item.module.scss';

interface MeetingTimelineItemProps {
  meeting: Meeting;
  onUpdateMeetingRecording: (id: string | null) => void;
}

export const MeetingRecording = ({
  meeting,
  onUpdateMeetingRecording,
}: MeetingTimelineItemProps): JSX.Element => {
  const uploadInputRef = React.useRef<HTMLInputElement>(null);
  const [summaryOpen, setSummaryOpen] = useState(false);

  const { isDraggingOver, handleDrag, handleDrop, handleInputFileChange } =
    useFileUpload({
      prevFiles: [],
      onBeginFileUpload: (data) => console.log('onBeginFileUpload', data),
      onFileUpload: (data) => {console.log("Upload done!"); onUpdateMeetingRecording(data.id);},
      onFileUploadError: () =>
        toast.error(
          'Something went wrong while uploading recording of a meeting',
        ),
      onFileRemove: () => onUpdateMeetingRecording(null),
      uploadInputRef,
    });

  const parseSummaryContent = (content?: string) => {
    if (!content) {
      return null;
    }
    let response;
    try {
      response = JSON.parse(content);
    } catch (e) {
      response = null;
    }
    return response;
  };
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

        {!!meeting?.describedBy.length && (
          <div className={styles.summaryItems}>
            {meeting?.describedBy.map((data, index) => {
              if (data.contentType === 'text/plain') {
                return (
                  <React.Fragment
                    key={`meeting-summary-item-${index}-${data.id}`}
                  >
                    <p>Summary:</p>
                    <div>{data.content}</div>
                  </React.Fragment>
                );
              }
              if (data.contentType && 'application/x-openline-action_items') {
                // @ts-expect-error fixme
                const actions = parseSummaryContent(data?.content);
                if (!actions) {
                  return (
                    <div key={`meeting-summary-content-unavailable-${data.id}`}>
                      Summary content unavailable
                    </div>
                  );
                }
                return (
                  <React.Fragment
                    key={`meeting-summary-content-action-items-${data.id}`}
                  >
                    <p>Action items:</p>
                    <ul>
                      {actions.action_list.map(
                        (action: string, index: number) => (
                          <li key={`meeting-analysis-action-item-${index}`}>
                            <Checkbox
                              type='checkbox'
                              label={action}
                              disabled
                              // @ts-expect-error fixme
                              onChange={() => null}
                            />
                          </li>
                        ),
                      )}
                    </ul>
                  </React.Fragment>
                );
              }
              return null;
            })}
          </div>
        )}
      </article>

      <div className={styles.collapsibleSection}>
        <div
          className={classNames(styles.transcriptionSection, {
            [styles.recordingUploaded]: meeting.recording,
            [styles.isDraggingOver]: isDraggingOver,
            [styles.collapsibleSectionWithSummary]: summaryOpen,
          })}
        ></div>

        <div className={styles.collapseExpandButtonWrapper}>
          <IconButton
            className={styles.collapseExpandButton}
            isSquare
            disabled={!meeting.recording}
            mode='secondary'
            size='xxxxs'
            icon={<ChevronDown width={24} height={24} />}
            onClick={() => setSummaryOpen(!summaryOpen)}
          />
        </div>
      </div>
      <section
        className={classNames(styles.summarySection, {
          [styles.summaryOpen]: summaryOpen,
        })}
      >
        {summaryOpen &&
          meeting.events.map((e, index) => {
            if (e.contentType === 'x-openline-transcript-element') {
              // @ts-expect-error fixme
              const transcript = parseSummaryContent(e.content);

              if (!transcript) {
                return (
                  <div key={`meeting-transcription-item-unavailable-${index}`}>
                    Transcription content could not be parsed
                  </div>
                );
              }

              return (
                <Message
                  key={`message-item-transcript-message-${index}`}
                  transcriptElement={{
                    ...transcript,
                    file_id: e?.includes?.[0].id,
                  }}
                  index={index}
                  contentType={e.contentType}
                  isLeft={false}
                />
              );
            }
            return null;
          })}
      </section>
    </>
  );
};
