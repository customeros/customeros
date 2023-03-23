import React, { useEffect, useRef, useState } from 'react';
import {
  ArrowLeft,
  ArrowRight,
  ChevronDown,
  ChevronUp,
  Message,
  Tooltip,
  VoiceWave,
} from '../../atoms';
import axios from 'axios';
import { toast } from 'react-toastify';
import { TimelineItem } from '../../atoms/timeline-item';
import { Skeleton } from '../../atoms/skeleton';
import { Content, ConversationItem, Props } from './types';
import styles from './conversation-timeline-item.module.scss';
import { AnalysisContent } from '../../atoms/message/AnalysisContent';
import classNames from 'classnames';
import { CallParties } from './CallParties';

export const PhoneConversationTimelineItem: React.FC<Props> = ({
  feedId,
  createdAt,
  feedInitiator,
}) => {
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const summaryRef = useRef<HTMLDivElement>(null);

  const [messages, setMessages] = useState([] as ConversationItem[]);
  const [summary, setSummary] = useState<{ analysis: Content } | undefined>();
  const [summaryExpanded, setSummaryExpanded] = useState(false);

  const [loadingMessages, setLoadingMessages] = useState(false);

  useEffect(() => {
    setLoadingMessages(true);
    axios
      .get(`/oasis-api/feed/${feedId}/item`)
      .then((res) => {
        const data = getSortedItems(res.data);
        setMessages(data ?? []);
        setLoadingMessages(false);
        (data ?? []).forEach((e) => {
          const decodedContent = decodeContent(e.content);
          if (decodedContent?.analysis?.type === 'summary') {
            setSummary(decodedContent);
          }
        });
      })
      .catch(() => {
        setLoadingMessages(false);
        toast.error('Something went wrong while loading feed item', {
          toastId: `conversation-timeline-item-feed-${feedId}`,
        });
      });
  }, [feedId]);

  //when a new message appears, scroll to the end of container
  useEffect(() => {
    if (messages) {
      setLoadingMessages(false);
    }
  }, [messages]);

  const timeFromLastTimestamp = new Date(1970, 0, 1).setSeconds(
    feedInitiator.lastTimestamp?.seconds,
  );

  const decodeContent = (content: string) => {
    let response;
    try {
      response = JSON.parse(content);
    } catch (e) {
      response = {
        dialog: {
          type: 'MESSAGE',
          mimetype: 'text/plain',
          body: content,
        },
      };
    }
    return response;
  };

  const handleToggleExpanded = () => {
    setSummaryExpanded(!summaryExpanded);
    if (summaryRef?.current && summaryExpanded) {
      summaryRef?.current?.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const getSortedItems = (data: Array<any>): Array<ConversationItem> => {
    return data
      .sort((a, b) => {
        const date1 =
          new Date(1970, 0, 1).setSeconds(a?.time?.seconds) ||
          timeFromLastTimestamp;
        const date2 =
          new Date(1970, 0, 1).setSeconds(b?.time?.seconds) ||
          timeFromLastTimestamp;
        return date2 - date1;
      })
      .filter((msg) => msg.type !== 1);
  };
  const left = messages.find((e) => e.direction === 0);
  const right = messages.find((e) => e.direction === 1);

  return (
    <div className='flex flex-column h-full w-full'>
      <div className='flex-grow-1 w-full'>
        {loadingMessages && (
          <div className='flex flex-column mb-2'>
            <div className='mb-2 flex justify-content-end'>
              <Skeleton height='40px' width='50%' />
            </div>
            <div className='mb-2 flex justify-content-start'>
              <Skeleton height='50px' width='40%' />
            </div>
            <div className='flex justify-content-end mb-2'>
              <Skeleton height='45px' width='50%' />
            </div>
            <div className='flex justify-content-start'>
              <Skeleton height='40px' width='45%' />
            </div>
          </div>
        )}

        <TimelineItem first createdAt={createdAt || timeFromLastTimestamp}>
          <div
            className={classNames(styles.contentWrapper, {
              [styles.expanded]: summaryExpanded,
            })}
          >
            <div className='flex flex-column w-full'>
              <div className={styles.summary} ref={summaryRef}>
                <div
                  className={classNames(styles.left, {
                    [styles.initiator]: messages?.[0]?.direction === 0,
                  })}
                >
                  <div className={styles.callPartyData}>
                    <CallParties
                      direction={0}
                      sender={left}
                      mode='PHONE_CALL'
                    />
                    <div className={styles.iconsWrapper}>
                      {messages[0]?.direction === 0 && (
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
                    [styles.initiator]: messages[0]?.direction === 1,
                  })}
                >
                  <div className={styles.callPartyData}>
                    <div className={styles.iconsWrapper}>
                      {messages[0]?.direction === 1 && (
                        <>
                          <ArrowLeft />
                          <VoiceWave />
                        </>
                      )}
                    </div>
                    <CallParties
                      direction={1}
                      sender={right}
                      mode='PHONE_CALL'
                    />
                  </div>
                </div>
              </div>
            </div>
            <Tooltip
              content={summary?.analysis?.body || ''}
              target={`#phone-summary-${feedId}`}
              position='bottom'
              showDelay={300}
              autoHide={false}
            />
            <button
              id={`phone-summary-${feedId}`}
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
              {summary && summary?.analysis && (
                <span>
                  Summary: <AnalysisContent analysis={summary?.analysis} />
                </span>
              )}
            </button>

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
                {!loadingMessages &&
                  getSortedItems(messages).map(
                    (msg: ConversationItem, index: number) => {
                      const lines = msg?.content.split('\n');

                      const filtered: string[] = lines.filter(
                        (line: string) => {
                          return line.indexOf('>') !== 0;
                        },
                      );
                      msg.content = filtered.join('\n').trim();

                      const time = new Date(1970, 0, 1).setSeconds(
                        msg?.time?.seconds,
                      );

                      return (
                        <Message
                          key={msg.id}
                          message={msg}
                          mode='PHONE_CALL'
                          date={time}
                          index={index}
                        />
                      );
                    },
                  )}
              </div>
            </section>
          </div>
        </TimelineItem>
      </div>
    </div>
  );
};
