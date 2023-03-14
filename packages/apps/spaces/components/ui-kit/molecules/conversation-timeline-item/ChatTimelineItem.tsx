import React, { useEffect, useState } from 'react';
import { Message } from '../../atoms';
import { TimelineItem } from '../../atoms/timeline-item';
import { Skeleton } from '../../atoms/skeleton';
import axios from 'axios';
import { toast } from 'react-toastify';
import { ConversationItem, Props } from './types';

export const ChatTimelineItem: React.FC<Props> = ({
  feedId,
  createdAt,
  first,
  feedInitiator,
}) => {
  const [messages, setMessages] = useState([] as ConversationItem[]);
  const [loadingMessages, setLoadingMessages] = useState(false);

  const makeSender = (msg: ConversationItem) => {
    return {
      loaded: true,
      email: msg.senderUsername.type == 0 ? msg.senderUsername.identifier : '',
      firstName: '',
      lastName: '',
      phoneNumber:
        msg.senderUsername.type == 1 ? msg.senderUsername.identifier : '',
    };
  };

  useEffect(() => {
    setLoadingMessages(true);
    axios
      .get(`/oasis-api/feed/${feedId}/item`)
      .then((res) => {
        setMessages(res.data ?? []);
        setLoadingMessages(false);
      })
      .catch(() => {
        setLoadingMessages(false);
        toast.error('Something went wrong while loading chat item', {
          toastId: `chat-conversation-timeline-item-feed-${feedId}`,
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

  return (
    <div className='flex flex-column h-full w-full'>
      <div className='flex-grow-1 w-full'>
        {loadingMessages && (
          <div className='flex flex-column mb-2'>
            <div className='mb-2 flex justify-content-end'>
              <Skeleton height='40px' width='50%' />
            </div>
          </div>
        )}

        {!loadingMessages &&
          getSortedItems(messages)
            .filter((msg) => msg.type !== 1)
            .map((msg: ConversationItem, index: number) => {
              const lines = msg?.content.split('\n');

              const filtered: string[] = lines.filter((line: string) => {
                return line.indexOf('>') !== 0;
              });
              msg.content = filtered.join('\n').trim();

              const time = new Date(1970, 0, 1).setSeconds(msg?.time?.seconds);

              const fl =
                first && (index === 0 || index === messages.length - 1);

              return (
                <TimelineItem
                  first={fl}
                  createdAt={createdAt || timeFromLastTimestamp}
                  key={msg.id}
                >
                  <Message
                    key={msg.id}
                    message={msg}
                    sender={makeSender(msg)}
                    date={time}
                    previousMessage={messages?.[index - 1]?.direction || null}
                    index={index}
                  />
                </TimelineItem>
              );
            })}
      </div>
    </div>
  );
};
