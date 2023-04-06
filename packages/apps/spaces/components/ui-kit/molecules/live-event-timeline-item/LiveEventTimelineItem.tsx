import React, { useEffect, useState } from 'react';
import { EmailTimelineItem } from '../email-timeline-item';
import { TimelineItem } from '../../atoms/timeline-item';
import useWebSocket from 'react-use-websocket';

interface Props {
  contactId?: string;
  source: string;
  first: boolean;
}

export type Time = {
  seconds: number;
};

export const LiveEventTimelineItem: React.FC<Props> = ({
  contactId,
  first,
}) => {
  const { lastMessage } = useWebSocket(
    `${process.env.NEXT_PUBLIC_WEBSOCKET_PATH}ws-participant/${contactId}`,
    {
      onOpen: () => console.log('Live events opened for contact ' + contactId),
      //Will attempt to reconnect on all close events, such as server shutting down
      shouldReconnect: (closeEvent) => true,
    },
  );

  const [liveEvents, setLiveEvents] = useState([] as any);
  const handleWebsocketLiveEvent = function (event: any) {
    console.log('Live events got new event:' + JSON.stringify(event));
    setLiveEvents((eventsList: any) => [...eventsList, event]);
  };

  useEffect(() => {
    if (
      lastMessage &&
      Object.keys(lastMessage).length !== 0 &&
      lastMessage.data.length > 0
    ) {
      //console.log('🏷️ ----- lastMessage: ', lastMessage?.data);
      handleWebsocketLiveEvent(JSON.parse(lastMessage?.data));
    }
  }, [lastMessage]);

  const getSortedItems = (data: Array<any>): Array<any> => {
    return data.sort((a, b) => {
      const date1 = new Date(1970, 0, 1).setSeconds(a?.createdAt?.seconds);
      const date2 = new Date(1970, 0, 1).setSeconds(b?.createdAt?.seconds);
      return date2 - date1;
    });
  };
  return (
    <div className='flex flex-column'>
      {
        // email
        getSortedItems(liveEvents).map((event: any, index: number) => {
          switch (event.channel) {
            case 'EMAIL': {
              const fl =
                first && (index === 0 || index === liveEvents.length - 1);
              return (
                <TimelineItem
                  first={fl}
                  createdAt={event.createdAt}
                  key={event.id}
                >
                  <EmailTimelineItem
                    content={event.content as string}
                    contentType={event.contentType as string}
                    sentBy={event.sentBy}
                    sentTo={event.sentTo}
                    interactionSession={{
                      name: event.interactionSession?.name,
                    }}
                  />
                </TimelineItem>
              );
            }
            case 'Voice': {
              // todo
              break;
            }
            default:
              return event.channel ? (
                <div>
                  Sorry, looks like &apos;{event.channel}&apos; live activity
                  type is not supported yet{' '}
                </div>
              ) : (
                ''
              );
          }
        })
      }
    </div>
  );
};
