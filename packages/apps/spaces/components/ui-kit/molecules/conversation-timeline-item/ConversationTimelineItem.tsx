import React, { useEffect, useRef, useState } from 'react';
import {
  ArrowLeft,
  ArrowRight,
  ChevronDown,
  ChevronUp,
  MessageIcon,
  Phone,
  Tooltip,
  VoiceWave,
} from '../../atoms';
import { TimelineItem } from '../../atoms/timeline-item';
import styles from './conversation-timeline-item.module.scss';
import { AnalysisContent } from '../../atoms/message/AnalysisContent';
import classNames from 'classnames';
import { TranscriptContent } from '../../atoms/message/TranscriptContent';
import {
  ConversationPartyEmail,
  ConversationPartyPhone,
} from './ConversationParty';

interface Content {
  dialog: {
    type?: string;
    mimetype: string;
    body: string;
  };
}

interface TranscriptElement {
  party: any;
  text: string;
  file_id?: string;
}

interface TranscriptV2 {
  transcript: Array<TranscriptElement>;
  file_id?: string;
}

interface Props {
  createdAt: string;
  content: Content | undefined;
  transcript: Array<TranscriptElement>;
  type: string;
  mode: 'PHONE_CALL' | 'CHAT';
  id: string;
  contentType?: string;
}

interface DataStateI {
  firstSendIndex: null | number;
  firstReceivedIndex: null | number;
  initiator: 'left' | 'right';
}

const getTranscript = (transcript:(TranscriptV2|Array<any>), contentType:string|undefined) => {
  console.log("getTranscript:", typeof transcript, "contentType: " + contentType);
  if (contentType === 'application/x-openline-transcript') {
    return transcript;
  } else if (contentType === 'application/x-openline-transcript-v2') {
    var transcript2 = transcript as TranscriptV2;
    return transcript2.transcript;
  }
}

export const ConversationTimelineItem: React.FC<Props> = ({
  createdAt,
  content,
  transcript = [],
  type,
  mode = 'PHONE_CALL',
  id,
  contentType,
}) => {
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const summaryRef = useRef<HTMLDivElement>(null);

  const [data, setData] = useState<DataStateI>({
    firstSendIndex: null,
    firstReceivedIndex: null,
    initiator: 'left',
  });
  const [summaryExpanded, setSummaryExpanded] = useState(false);
  const handleToggleExpanded = () => {
    setSummaryExpanded(!summaryExpanded);
    if (summaryRef?.current && summaryExpanded) {
      summaryRef?.current?.scrollIntoView({ behavior: 'smooth' });
    }
  };

  useEffect(() => {
    if (data.firstSendIndex === null) {
      const left = getTranscript(transcript, contentType).findIndex(
        (e: TranscriptElement) => e?.party?.tel,
      );
      const right = getTranscript(transcript, contentType).findIndex(
        (e: TranscriptElement) => e?.party?.mailto,
      );

      setData({
        firstSendIndex: left,
        firstReceivedIndex: right,
        initiator: left === 0 ? 'left' : 'right',
      });
    }
  }, []);
  // fixme for some reason it does not work whe put in state
  console.log("transcript", transcript);
  const left = getTranscript(transcript, contentType)?.find((e: TranscriptElement) => e?.party?.tel);
  const right = getTranscript(transcript, contentType)?.find((e: TranscriptElement) => e?.party?.mailto);
  //const right=false, left = false;
  return (
    <div className='flex flex-column w-full'>
      <TimelineItem first createdAt={createdAt}>
        <div
          className={classNames(styles.contentWrapper, {
            [styles.expanded]: summaryExpanded,
          })}
        >
          {type === 'summary' && (
            <>
              <div className='flex flex-column w-full'>
                <div className={styles.summary} ref={summaryRef}>
                  <div
                    className={classNames(styles.left, {
                      [styles.initiator]: data.initiator === 'left',
                    })}
                  >
                    <div className={styles.callPartyData}>
                      <ConversationPartyPhone tel={left?.party.tel} />

                      <div className={styles.iconsWrapper}>
                        {transcript?.[0]?.party.tel && (
                          <>
                            <VoiceWave />
                            <ArrowRight />
                          </>
                        )}
                      </div>
                    </div>
                  </div>

                  <div
                    className={classNames(styles.right, {
                      [styles.initiator]: data.initiator === 'right',
                    })}
                  >
                    <div className={styles.callPartyData}>
                      <div className={styles.iconsWrapper}>
                        {!transcript?.[0]?.party.tel && (
                          <>
                            <ArrowLeft />
                            <VoiceWave />
                          </>
                        )}
                      </div>
                      <ConversationPartyEmail
                        email={(right?.party.mailto || '').toLowerCase()}
                      />
                    </div>
                  </div>
                </div>
              </div>

              {content && (
                <Tooltip
                  content={content.dialog?.body || ''}
                  target={`#phone-summary-${id}`}
                  position='bottom'
                  showDelay={300}
                  autoHide={false}
                />
              )}

              <button
                id={`phone-summary-${id}`}
                className={styles.folderTab}
                role='button'
                onClick={handleToggleExpanded}
              >
                {summaryExpanded ? (
                  <ChevronUp
                    style={{
                      color: '#3A8745',
                      minWidth: '23px',
                      transform: 'scale(0.8)',
                    }}
                  />
                ) : (
                  <ChevronDown
                    style={{
                      color: '#3A8745',
                      minWidth: '23px',
                      transform: 'scale(0.8)',
                    }}
                  />
                )}

                {content && (
                  <span>
                    Summary: <AnalysisContent analysis={content.dialog} />
                  </span>
                )}
              </button>
            </>
          )}

          <section
            ref={messagesContainerRef}
            className={classNames(styles.transcriptionContainer, {
              [styles.transcriptionContainerOpen]: summaryExpanded,
            })}
            style={{
              maxHeight: summaryExpanded
                ? `${messagesContainerRef?.current?.scrollHeight}px`
                : 0,
            }}
          >
            <div className={styles.messages}>
              <TranscriptContent
                contentType={contentType}
                messages={getTranscript(transcript, contentType)}
                firstIndex={{
                  received: 0,
                  send: 1,
                }}
              >
                {mode === 'CHAT' ? <MessageIcon /> : <Phone />}
              </TranscriptContent>
            </div>
          </section>
        </div>
      </TimelineItem>
    </div>
  );
};
