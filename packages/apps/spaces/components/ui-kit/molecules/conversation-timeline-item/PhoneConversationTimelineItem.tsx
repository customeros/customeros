import React, { useEffect, useState } from 'react';
import {
  ArrowLeft,
  ArrowRight,
  Avatar,
  ChevronDown,
  ChevronUp,
  IconButton,
  Message,
  VoiceWave,
} from '../../atoms';
import axios from 'axios';
import { toast } from 'react-toastify';
import { TimelineItem } from '../../atoms/timeline-item';
import { Skeleton } from '../../atoms/skeleton';
import { ConversationItem, Props } from './types';
import styles from './conversation-timeline-item.module.scss';
import { AnalysisContent } from '../../atoms/message/AnalysisContent';
import classNames from 'classnames';
import { CallParties } from './CallParties';
export const PhoneConversationTimelineItem: React.FC<Props> = ({
  feedId,
  createdAt,
  first,
  feedInitiator,
}) => {
  const [messages, setMessages] = useState([] as ConversationItem[]);
  const [summary, setSummary] = useState();
  const [summaryExpanded, setSummaryExpanded] = useState(false);

  const [loadingMessages, setLoadingMessages] = useState(false);

  useEffect(() => {
    setLoadingMessages(true);
    axios
      .get(`/oasis-api/feed/${feedId}/item`)
      .then((res) => {
        setMessages(res.data ?? []);
        setLoadingMessages(false);
        (res.data ?? []).forEach((e) => {
          const z = decodeContent(e.content);
          if (z?.analysis?.type === 'summary') {
            setSummary(z);
          }
        });
      })
      .catch(() => {
        setLoadingMessages(false);
        toast.error('Something went wrong while loading feed item', {
          toastId: `conversation-timeline-item-feed-${feedId}`,
        });
      });
  }, []);

  //when a new message appears, scroll to the end of container
  useEffect(() => {
    if (messages) {
      setLoadingMessages(false);
    }
  }, [messages]);

  const timeFromLastTimestamp = new Date(1970, 0, 1).setSeconds(
    feedInitiator.lastTimestamp?.seconds,
  );
  // const timeFromFirstTimestamp = new Date(1970, 0, 1).setSeconds(
  //   feedInitiator.lastTimestamp?.seconds,
  // );
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

  const getSortedItems = (data: Array<any>): Array<ConversationItem> => {
    return data.sort((a, b) => {
      const date1 =
        new Date(1970, 0, 1).setSeconds(a?.time?.seconds) ||
        timeFromLastTimestamp;
      const date2 =
        new Date(1970, 0, 1).setSeconds(b?.time?.seconds) ||
        timeFromLastTimestamp;
      return date2 - date1;
    });
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
          <div className={styles.contentWrapper}>
            <div className='flex flex-column w-full'>
              <div className={styles.summary}>
                <div className={styles.left}>
                  <div className={styles.callPartyData}>
                    <CallParties direction={left?.direction} sender={left} />
                    <VoiceWave />
                    <ArrowRight />
                  </div>
                </div>

                <div className={styles.right}>
                  <div className={styles.callPartyData}>
                    <div></div>
                    <CallParties direction={right?.direction} sender={right} />
                  </div>
                </div>
              </div>
            </div>
            <div className={styles.folderTab}>
              <IconButton
                mode='text'
                size='xxxs'
                onClick={() => setSummaryExpanded(!summaryExpanded)}
                style={{ color: '#3A8745' }}
                icon={summaryExpanded ? <ChevronUp /> : <ChevronDown />}
              />
              {summary && summary?.analysis && (
                <span>
                  Summary: <AnalysisContent analysis={summary?.analysis} />
                </span>
              )}
            </div>

            <section
              className={classNames(styles.transcriptionContainer, {
                [styles.transcriptionContainerOpen]: summaryExpanded,
              })}
            >
              {!loadingMessages &&
                getSortedItems(messages)
                  .filter((msg) => msg.type !== 1)
                  .map((msg: ConversationItem, index: number) => {
                    const lines = msg?.content.split('\n');

                    const filtered: string[] = lines.filter((line: string) => {
                      return line.indexOf('>') !== 0;
                    });
                    msg.content = filtered.join('\n').trim();

                    const time = new Date(1970, 0, 1).setSeconds(
                      msg?.time?.seconds,
                    );

                    const fl =
                      first && (index === 0 || index === messages.length - 1);
                    const sender = {
                      email:
                        msg.senderUsername.type == 0
                          ? msg.senderUsername.identifier
                          : '',
                      firstName: '',
                      lastName: '',
                      phoneNumber:
                        msg.senderUsername.type == 1
                          ? msg.senderUsername.identifier
                          : '',
                    };
                    const initiatior = {
                      email:
                        msg.initiatorUsername.type == 0
                          ? msg.initiatorUsername.identifier
                          : '',
                      firstName: '',
                      lastName: '',
                      phoneNumber:
                        msg.initiatorUsername.type == 1
                          ? msg.initiatorUsername.identifier
                          : '',
                    };

                    // console.log('🏷️ ----- sender: ', sender);
                    // console.log('🏷️ ----- message: ', msg);

                    return (
                      <div>
                        <Message
                          key={msg.id}
                          message={msg}
                          sender={sender}
                          date={time}
                          previousMessage={
                            messages?.[index - 1]?.direction || null
                          }
                          index={index}
                        />
                      </div>
                    );
                  })}
            </section>
          </div>
        </TimelineItem>
      </div>
    </div>
  );
};
