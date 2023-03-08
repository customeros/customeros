import React, { useEffect, useState } from 'react';
import { Button, Message, Reply } from '../../atoms';
import axios from 'axios';
import { gql } from 'graphql-request';
import { toast } from 'react-toastify';
import { EmailTimelineItem } from '../email-timeline-item';
import { TimelineItem } from '../../atoms/timeline-item';
import { Skeleton } from '../../atoms/skeleton';
import { useRecoilState, useRecoilValue } from 'recoil';
import {
  editorEmail,
  editorMode,
  EditorMode,
  userData,
} from '../../../../state';

interface Props {
  feedId: string;
  source: string;
  first: boolean;
  createdAt: any;
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

export type FeedPostRequest = {
  username: string;
  message: string;
  channel: string;
  direction: string;
  destination: string[];
  replyTo?: string;
};
export const ConversationTimelineItem: React.FC<Props> = ({
  feedId,
  source,
  createdAt,
  first,
}) => {
  const [editorModeState, setEditorMode] = useRecoilState(editorMode);
  const [emailEditorData, setEmailEditorData] = useRecoilState(editorEmail);
  const loggedInUserData = useRecoilValue(userData);
  const [feedInitiator, setFeedInitiator] = useState<any>({
    loaded: false,
    email: '',
    firstName: '',
    lastName: '',
    phoneNumber: '',
    lastTimestamp: null,
  });

  const [messages, setMessages] = useState([] as ConversationItem[]);
  const [participants, setParticipants] = useState([] as Participant[]);
  const [isSending, setSending] = useState(false);

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
      .get(`/oasis-api/feed/${feedId}`)
      .then((res) => {
        const feedItem = res.data as FeedItem;

        if (feedItem.initiatorType !== 'CONTACT') {
          setFeedInitiator({
            loaded: true,
            email:
              feedItem.initiatorUsername.type == 0
                ? feedItem.initiatorUsername.identifier
                : '',
            firstName: feedItem.initiatorFirstName,
            lastName: feedItem.initiatorLastName,
            phoneNumber:
              feedItem.initiatorUsername.type == 1
                ? feedItem.initiatorUsername.identifier
                : '',
            lastTimestamp: feedItem.lastTimestamp,
          });
        }

        if (feedItem.initiatorType === 'CONTACT') {
          const query = gql`
            query GetContactDetails($email: String!) {
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
            }
          `;

          //TODO
          setFeedInitiator({
            loaded: true,
            firstName: 'EDIT',
            lastName: 'EDIT',
            email: undefined,
            phoneNumber: undefined,
            lastTimestamp: feedItem.lastTimestamp,
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
      })
      .catch((reason: any) => {
        //todo log on backend

        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });

    axios
      .get(`/oasis-api/feed/${feedId}/participants`)
      .then((res) => {
        if (res.data && res.data.participants) {
          const list: Participant[] = [];
          for (const participant of res.data.participants) {
            if (participant.email !== loggedInUserData.identity) {
              //@ts-expect-error this will be fixed after switching to new paginated timeline
              list.push({ email: participant.email });
            }
          }

          setParticipants(list);
        }
      })
      .catch((reason: any) => {
        //todo log on backend
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
    axios
      .get(`/oasis-api/feed/${feedId}/item`)
      .then((res) => {
        setMessages(res.data ?? []);
        setLoadingMessages(false);
      })
      .catch((reason: any) => {
        setLoadingMessages(false);
        toast.error('Something went wrong while loading feed item', {
          toastId: `conversation-timeline-item-feed-${feedId}`,
        });
      });
  }, []);

  const handleSendMessage = (
    text: string,
    onSuccess: () => void,
    destination = [],
    replyTo: null | string,
  ) => {
    if (!text) return;
    const message: FeedPostRequest = {
      channel: 'EMAIL',
      username: loggedInUserData.identity,
      message: text,
      direction: 'OUTBOUND',
      destination: destination,
    };
    console.log('🏷️ ----- replyTo: ', replyTo);
    if (replyTo) {
      message.replyTo = replyTo;
    }

    setSending(true);

    axios
      .post(`/oasis-api/feed/${feedId}/item`, message)
      .then((res) => {
        console.log(res);
        if (res.data) {
          setMessages((messageList: any) => [...messageList, res.data]);
          onSuccess();
          setEditorMode({
            submitButtonLabel: 'Log into timeline',
            mode: EditorMode.Note,
          });
          setEmailEditorData({ ...emailEditorData, to: [], subject: '' });
          toast.success('Email sent!');
          setSending(false);
        }
      })
      .catch((reason) => {
        //todo log on backend
        setSending(false);
        toast.error('Something went wrong while sending email');
      });
  };

  //when a new message appears, scroll to the end of container
  useEffect(() => {
    if (messages && feedInitiator.loaded) {
      setLoadingMessages(false);
    }
  }, [messages, feedInitiator]);

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

        <div className='flex flex-column'>
          {
            // email
            !loadingMessages &&
              getSortedItems(messages)
                .filter((msg: ConversationItem) => msg.type === 1)
                .map((msg: ConversationItem, index: number) => {
                  const emailData = JSON.parse(msg.content);

                  const date =
                    new Date(1970, 0, 1).setSeconds(msg?.time?.seconds) ||
                    timeFromLastTimestamp;
                  const fl =
                    first && (index === 0 || index === messages.length - 1);
                  return (
                    <TimelineItem
                      first={fl}
                      createdAt={date}
                      //@ts-expect-error fixme later
                      key={msg?.messageId?.conversationEventId}
                    >
                      <EmailTimelineItem
                        emailContent={emailData.html}
                        sender={emailData.from || 'Unknown'}
                        recipients={emailData.to}
                        cc={emailData?.cc}
                        bcc={emailData?.bcc}
                        subject={emailData.subject}
                      >
                        <Button
                          mode='link'
                          onClick={() => {
                            // TODO add cc and bcc

                            setEmailEditorData({
                              //@ts-expect-error fixme later
                              handleSubmit: handleSendMessage,
                              to: [emailData.from],
                              subject: emailData.subject,
                              respondTo:
                                //@ts-expect-error fixme later
                                msg?.messageId?.conversationEventId || null,
                            });
                            setEditorMode({
                              mode: EditorMode.Email,
                              submitButtonLabel: 'Send',
                            });
                          }}
                        >
                          Respond
                        </Button>
                      </EmailTimelineItem>
                    </TimelineItem>
                  );
                })
          }
        </div>

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
