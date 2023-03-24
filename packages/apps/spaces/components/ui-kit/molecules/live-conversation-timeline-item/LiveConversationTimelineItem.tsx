import React, { useEffect, useState } from 'react';
import { MessageDeprecate } from '../../atoms';
import { Skeleton } from 'primereact/skeleton';
import { TimelineItem } from '../../atoms/timeline-item';
import useWebSocket from 'react-use-websocket';
import { EmailTimelineItemToDeprecate } from '../email-timeline-item-to-deprecate';

interface Props {
  contactId?: string;
  source: string;
  first: boolean;
}

export type Time = {
  seconds: number;
};

export type Participant = {
  type: number;
  identifier: string;
};

export type ConversationItem = {
  id: string;
  conversationId: string;
  type: number;
  subtype: number;
  content: string;
  direction: number;
  time: Time;
  senderType: number;
  senderId: string;
  senderUsername: Participant;
};

export type FeedItem = {
  id: string;
  initiatorFirstName: string;
  initiatorLastName: string;
  initiatorUsername: Participant;
  initiatorType: string;
  lastSenderFirstName: string;
  lastSenderLastName: string;
  lastContentPreview: string;
  lastTimestamp: Time;
};

export const LiveConversationTimelineItem: React.FC<Props> = ({
  contactId,
  source,
  first,
}) => {
  const { lastMessage } = useWebSocket(
    `${process.env.NEXT_PUBLIC_WEBSOCKET_PATH}ws-participant/${contactId}`,
    {
      onOpen: () => console.log('Websocket opened'),
      //Will attempt to reconnect on all close events, such as server shutting down
      shouldReconnect: (closeEvent) => true,
    },
  );

  const [messages, setMessages] = useState([] as ConversationItem[]);

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
  const handleWebsocketMessage = function (msg: any) {
    console.log('Got new message:' + JSON.stringify(msg));
    setMessages((messageList: any) => [...messageList, msg]);
  };

  useEffect(() => {
    if (
      lastMessage &&
      Object.keys(lastMessage).length !== 0 &&
      lastMessage.data.length > 0
    ) {
      // console.log('🏷️ ----- lastMessage: ', lastMessage);
      handleWebsocketMessage(JSON.parse(lastMessage?.data));
    }
  }, [lastMessage]);

  const getSortedItems = (data: Array<any>): Array<ConversationItem> => {
    return data.sort((a, b) => {
      const date1 = new Date(1970, 0, 1).setSeconds(a?.time?.seconds);
      const date2 = new Date(1970, 0, 1).setSeconds(b?.time?.seconds);
      return date2 - date1;
    });
  };
  return (
    <div className='flex flex-column w-full'>
      <div className='flex-grow-1 w-full'>
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

        <div className='flex flex-column'>
          {
            // email
            getSortedItems(messages)
              .filter((msg: ConversationItem) => msg.type === 1)
              .map((msg: ConversationItem, index: number) => {
                const emailData = JSON.parse(msg.content);
                const date = new Date(1970, 0, 1).setSeconds(
                  msg?.time?.seconds,
                );
                const fl =
                  first && (index === 0 || index === messages.length - 1);
                return (
                  <TimelineItem first={fl} createdAt={date} key={msg.id}>
                    {/*TODO switch to EmailTimelineItem when the backend migration is done*/}
                    <EmailTimelineItemToDeprecate
                      emailContent={emailData.html}
                      sender={emailData.from || 'Unknown'}
                      recipients={emailData.to}
                      cc={emailData?.cc}
                      bcc={emailData?.bcc}
                      subject={emailData.subject}
                    />
                  </TimelineItem>
                );
              })
          }
        </div>

        {getSortedItems(messages)
          .filter((msg) => msg.type !== 1)
          .map((msg: ConversationItem, index: number) => {
            const lines = msg?.content.split('\n');

            const filtered: string[] = lines.filter((line: string) => {
              return line.indexOf('>') !== 0;
            });
            msg.content = filtered.join('\n').trim();

            const time = new Date(1970, 0, 1).setSeconds(msg?.time?.seconds);

            const fl = first && (index === 0 || index === messages.length - 1);

            return (
              <TimelineItem first={fl} createdAt={time} key={msg.id}>
                <MessageDeprecate
                  key={msg.id}
                  message={msg}
                  date={time}
                  index={index}
                  mode='LIVE_CONVERSATION'
                />
                <span className='text-sm '>
                  {' '}
                  Source: {source?.toLowerCase() || 'unknown'}
                </span>
              </TimelineItem>
            );
          })}
      </div>
    </div>
  );
};
