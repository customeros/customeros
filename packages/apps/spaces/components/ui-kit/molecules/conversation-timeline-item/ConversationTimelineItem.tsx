import React, { useEffect, useState } from 'react';
import { Message } from '../../atoms';
import axios from 'axios';
// import { FeedItem } from '../../../models/feed-item';
import { gql } from 'graphql-request';
import { toast } from 'react-toastify';
import useWebSocket from 'react-use-websocket';
// import { ConversationItem } from '../../../models/conversation-item';
import { Skeleton } from 'primereact/skeleton';
// import { useGraphQLClient } from '../../../utils/graphQLClient';
import { EmailTimelineItem } from '../email-timeline-item';
import { TimelineItem } from '../../atoms/timeline-item';
interface Props {
  feedId: string;
  source: string;
  fistOrLast: boolean;
  createdAt: any;
}

export type Time = {
  seconds: number,
}

export type ConversationItem = {
  id:             string,
  conversationId: string,
  type:           number,
  subtype:        number,
  content:        string,
  direction:      number,
  time:           Time,
  senderType:     number,
  senderId:       string,
  senderUserName?:       string,
}

export type FeedItem = {
  id: string;
  initiatorFirstName: string;
  initiatorLastName: string;
  initiatorUsername: string;
  initiatorType: string;
  lastSenderFirstName: string;
  lastSenderLastName: string;
  lastContentPreview: string;
  lastTimestamp: Time;
}


export const ConversationTimelineItem: React.FC<Props> = (
    { feedId, source, createdAt, fistOrLast}) => {

  const [feedInitiator, setFeedInitiator] = useState<any>({
    loaded: false,
    email: '',
    firstName: '',
    lastName: '',
    phoneNumber: '',
    lastTimestamp: null
  });

  const [messages, setMessages] = useState([] as ConversationItem[]);

  const [loadingMessages, setLoadingMessages] = useState(false)

  useEffect(() => {
    setLoadingMessages(true);
    axios.get(`/oasis-api/feed/${feedId}`)
    .then(res => {
      const feedItem = res.data as FeedItem;

      if(feedItem.initiatorType !== 'CONTACT') {
        setFeedInitiator({
          loaded: true,
          email: feedItem.initiatorUsername,
          firstName: feedItem.initiatorFirstName,
          lastName: feedItem.initiatorLastName,
          phoneNumber: '',
          lastTimestamp: feedItem.lastTimestamp
        })
      }

      if (feedItem.initiatorType === 'CONTACT') {

        const query = gql`query GetContactDetails($email: String!) {
          contact_ByEmail(email: $email) {
            id
            firstName
            lastName
            emails {
              email
            }
            phoneNumbers {
              e164
            }
          }
        }`

        //TODO
            setFeedInitiator({
              loaded: true,
              firstName: "EDIT",
              lastName: "EDIT",
              email: undefined,
              phoneNumber:  undefined,
              lastTimestamp: feedItem.lastTimestamp
            });

        // client.request(query, {email: feedItem.initiatorUsername}).then((response: any) => {
        //   if (response.contact_ByEmail) {
        //     setFeedInitiator({
        //       loaded: true,
        //       firstName: response.contact_ByEmail.firstName,
        //       lastName: response.contact_ByEmail.lastName,
        //       email: response.contact_ByEmail.emails[0]?.email ?? undefined,
        //       phoneNumber: response.contact_ByEmail.phoneNumbers[0]?.e164 ?? undefined,
        //       lastTimestamp: feedItem.lastTimestamp
        //     });
        //   } else {
        //     //todo log on backend
        //     toast.error("There was a problem on our side and we are doing our best to solve it!");
        //   }
        // }).catch(reason => {
        //   //todo log on backend
        //   toast.error("There was a problem on our side and we are doing our best to solve it!");
        // });

        //TODO move initiator in index
      }

    }).catch((reason: any) => {
      //todo log on backend

      toast.error("There was a problem on our side and we are doing our best to solve it!");
    });

    axios.get(`/oasis-api/feed/${feedId}/item`)
        .then(res => {
          setMessages(res.data ?? []);
          setLoadingMessages(false)
        }).catch((reason: any) => {
      setLoadingMessages(false)
      toast.error("There was a problem on our side and we are doing our best to solve it!");
    });
  }, []);

  //when a new message appears, scroll to the end of container
  useEffect(() => {
    if (messages && feedInitiator.loaded) {

      setLoadingMessages(false);
    }
  }, [messages, feedInitiator]);

  const timeFromLastTimestamp = new Date(1970, 0, 1)
      .setSeconds(feedInitiator.lastTimestamp?.seconds);

  const getSortedItems = (data: Array<any>): Array<ConversationItem> => {
    return data.sort((a, b) => {
      const date1 =  new Date(1970, 0, 1)
          // @ts-ignore
          .setSeconds(a?.time?.seconds) || timeFromLastTimestamp;
      const date2 =  new Date(1970, 0, 1)
          // @ts-ignore
          .setSeconds(b?.time?.seconds) || timeFromLastTimestamp;
      return  date2  - date1;
    })
  }
  return (
      <div className='flex flex-column h-full w-full'>
        <div className="flex-grow-1 w-full">
          {
              loadingMessages &&
              <div className="flex flex-column mb-2">
                <div className="mb-2 flex justify-content-end">
                  <Skeleton height="40px" width="50%" />
                </div>
                <div className="mb-2 flex justify-content-start">
                  <Skeleton height="50px" width="40%" />
                </div>
                <div className="flex justify-content-end mb-2">
                  <Skeleton height="45px" width="50%" />
                </div>
                <div className="flex justify-content-start">
                  <Skeleton height="40px" width="45%" />
                </div>
              </div>
          }

          <div className="flex flex-column">
            {   // email
                !loadingMessages &&
                getSortedItems(messages).filter((msg:ConversationItem) => msg.type === 1).map((msg: ConversationItem, index: number) => {
                  const emailData = JSON.parse(msg.content)
                  const date =  new Date(1970, 0, 1)
                      .setSeconds(msg?.time?.seconds) || timeFromLastTimestamp;
                  const fl = fistOrLast && (index === 0 || index === messages.length - 1);
                  return (
                      <TimelineItem fistOrLast={fl}
                                    createdAt={date}
                                    style={{paddingBottom: '8px'}}
                                    key={msg.id}>
                        <EmailTimelineItem
                            emailContent={emailData.html}
                            sender={emailData.from || 'Unknown'}
                            recipients={emailData.to}
                            cc={emailData?.cc}
                            bcc={emailData?.bcc}
                            subject={emailData.subject}
                        />
                      </TimelineItem>
                  )
                })
            }
          </div>

          {
              !loadingMessages &&
              getSortedItems(messages).filter(msg => msg.type !== 1).map((msg: ConversationItem, index: number) => {
                const lines = msg?.content.split('\n');

                const filtered: string[] = lines.filter((line: string) => {
                  return line.indexOf('>') !== 0;
                });
                msg.content = filtered.join('\n').trim();

                const time = new Date(1970, 0, 1).setSeconds(msg?.time?.seconds);

                const fl = fistOrLast && (index === 0 || index === messages.length - 1);

                return (
                    <TimelineItem fistOrLast={fl} createdAt={createdAt || timeFromLastTimestamp} key={msg.id}>
                      <Message key={msg.id}
                               message={msg}
                               feedInitiator={feedInitiator}
                               date={time}
                               previousMessage={messages?.[index - 1]?.direction || null}
                               index={index} />
                      <span className="text-sm "> Source: {source?.toLowerCase() || 'unknown'}</span>
                    </TimelineItem>
                )
              })
          }
        </div>
      </div>
  );
};
