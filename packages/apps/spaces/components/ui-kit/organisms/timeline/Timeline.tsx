import React, { useRef } from 'react';
import { useStickyScroll } from '../../../../hooks/useStickyScroll';

import {
  ConversationTimelineItem,
  EmailTimelineItem,
  LiveConversationTimelineItem,
  NoteTimelineItem,
  PhoneCallTimelineItem,
  WebActionTimelineItem,
} from '../../molecules';
import { TimelineItem } from '../../atoms/timeline-item';
import { TicketTimelineItem } from '../../molecules/ticket-timeline-item';
import styles from './timeline.module.scss';
import { InteractionTimelineItem } from '../../molecules/interaction-timeline-item';
import { useInfiniteScroll } from './useInfiniteScroll';
import { Skeleton } from '../../atoms/skeleton';
import { TimelineStatus } from './timeline-status';
import classNames from 'classnames';
import { PhoneConversationTimelineItem } from '../../molecules/conversation-timeline-item/PhoneConversationTimelineItem';
import { EmailTimelineItemToDeprecate } from '../../molecules/email-timeline-item-to-deprecate';
import { EmailTimelineItemTemp } from '../../molecules/conversation-timeline-item/EmailTimelineItemTemp';

interface Props {
  loading: boolean;
  noActivity: boolean;
  id?: string;
  loggedActivities: Array<any>;
  notifyChange?: (id: any) => void;
  onLoadMore: (ref: any) => void;
  contactName?: string;
  mode: 'CONTACT' | 'ORGANIZATION';
}

export const Timeline = ({
  loading,
  noActivity,
  loggedActivities,
  id,
  onLoadMore,
  contactName = '',
  mode,
}: Props) => {
  const timelineContainerRef = useRef(null);
  const containerRef = useRef(null);
  const infiniteScrollElementRef = useRef(null);
  // @ts-expect-error revisit later
  useStickyScroll(containerRef, loggedActivities || []);
  useInfiniteScroll({
    element: infiniteScrollElementRef,
    isFetching: loading,
    callback: () => {
      if (loggedActivities.length > 10) {
        onLoadMore(containerRef);
      }
    },
  });

  const getTimelineItemByType = (type: string, data: any, index: number) => {
    switch (type) {
      case 'Note':
        return (
          <TimelineItem first={index == 0} createdAt={data?.createdAt}>
            <NoteTimelineItem
              noteContent={data.html}
              createdAt={data.createdAt}
              createdBy={data?.createdBy}
              id={data.id}
              source={data?.source}
              noted={data?.noted}
            />
          </TimelineItem>
        );
      case 'Conversation':
        // TODO move to interaction event once we have the data in backend
        // if (data.channel === 'WEB_CHAT') {
        //   return (
        //     <ChatTimelineItem
        //       first={index == 0}
        //       feedId={data.id}
        //       source={data.source}
        //       createdAt={data?.startedAt}
        //       feedInitiator={{
        //         firstName: data.initiatorFirstName,
        //         lastName: data.initiatorLastName,
        //         phoneNumber: data.initiatorUsername.identifier,
        //         lastTimestamp: data.lastTimestamp,
        //       }}
        //     />
        //   );
        // }
        if (data.channel === 'EMAIL') {
          return (
            <EmailTimelineItemTemp
              first={index == 0}
              feedId={data.id}
              source={data.source}
              createdAt={data?.startedAt}
              feedInitiator={{
                firstName: data.initiatorFirstName,
                lastName: data.initiatorLastName,
                phoneNumber: data.initiatorUsername.identifier,
                lastTimestamp: data.lastTimestamp,
              }}
            />
          );
        }
        // TODO move to interaction event once we have the data in backend
        if (data.channel === 'VOICE') {
          return (
            <PhoneConversationTimelineItem
              first={index == 0}
              feedId={data.id}
              source={data.source}
              createdAt={data?.startedAt}
              feedInitiator={{
                firstName: data.initiatorFirstName,
                lastName: data.initiatorLastName,
                phoneNumber: data.initiatorUsername.identifier,
                lastTimestamp: data.lastTimestamp,
              }}
            />
          );
        }
        return null;

      case 'Analysis': {
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
        if (data.analysisType === 'transcript') {
          return null;
        }

        const transcriptForSummary = loggedActivities
          .filter((e) => e.__typename === 'Analysis')
          .filter((e) => e.analysisType !== 'summary')
          .find((e) => e.describes[0].id === data.describes[0].id);

        if (!transcriptForSummary?.content) {
          return;
        }

        return (
          <ConversationTimelineItem
            id={data.id}
            content={decodeContent(data.content)}
            transcript={decodeContent(transcriptForSummary.content)}
            type={data.analysisType}
            createdAt={data?.createdAt}
            mode='PHONE_CALL' // fixme - mode will be assessed from data inside the component (on message base)
          />
        );
      }

      case 'LiveConversation':
        return (
          <LiveConversationTimelineItem
            first={index == 0}
            contactId={id}
            source={data.source}
          />
        );
      case 'PageView':
        return (
          <TimelineItem first={index == 0} createdAt={data?.startedAt}>
            <WebActionTimelineItem {...data} contactName={contactName} />
          </TimelineItem>
        );
      case 'InteractionSession':
        return (
          <TimelineItem first={index == 0} createdAt={data?.startedAt}>
            <InteractionTimelineItem
              {...data}
              contactId={contactName && id}
              organizationId={!contactName && id}
            />
          </TimelineItem>
        );
      case 'Ticket':
        return (
          <TimelineItem first={index == 0} createdAt={data?.createdAt}>
            <TicketTimelineItem {...data} />
          </TimelineItem>
        );

      case 'InteractionEvent':
        if (data.channel === 'EMAIL') {
          return (
            <TimelineItem first={index == 0} createdAt={data?.createdAt}>
              <EmailTimelineItem
                {...data}
                contactId={mode === 'CONTACT' && id}
              />
            </TimelineItem>
          );
        }
        if (data.channel === 'VOICE') {
          const from =
            data.sentBy && data.sentBy.length > 0
              ? data.sentBy
                  .map((p: any) => {
                    if (
                      p.__typename === 'EmailParticipant' &&
                      p.emailParticipant
                    ) {
                      return p.emailParticipant.email;
                    }
                    return '';
                  })
                  .join('; ')
              : '';

          const to =
            data.sentTo && data.sentTo.length > 0
              ? data.sentTo
                  .map((p: any) => {
                    if (
                      p.__typename === 'EmailParticipant' &&
                      p.emailParticipant
                    ) {
                      return p.emailParticipant.email;
                    } else if (
                      p.__typename === 'ContactParticipant' &&
                      p.contactParticipant
                    ) {
                      if (
                        p.contactParticipant.name &&
                        p.contactParticipant.name !== ''
                      ) {
                        return p.contactParticipant.name;
                      } else {
                        return (
                          p.contactParticipant.firstName +
                          ' ' +
                          p.contactParticipant.lastName
                        );
                      }
                    }
                    return '';
                  })
                  .join('; ')
              : '';

          //we are using this to render the phone calls manually created by the user
          return (
            <ConversationTimelineItem
              id={data.id}
              content={undefined}
              transcript={[
                {
                  text: data.content,
                  party: {
                    tel: from,
                    mailto: to,
                  },
                },
              ]}
              type={'summary'} //fixme: this is used to get the same style as the summary of  phone call
              createdAt={data?.createdAt}
              mode='PHONE_CALL' // fixme - mode will be assessed from data inside the component (on message base)
            />
          );
        } else {
          return null;
        }

      default:
        return type ? (
          <div>
            Sorry, looks like &apos;{type}&apos; activity type is not supported
            yet{' '}
          </div>
        ) : (
          ''
        );
    }
  };

  return (
    <div ref={timelineContainerRef} className={styles.timeline}>
      {!loading && noActivity && <TimelineStatus status='no-activity' />}
      <div
        className={classNames(styles.timelineContent, {
          [styles.scrollable]: !noActivity,
        })}
        ref={containerRef}
      >
        {!!loggedActivities.length && (
          <div
            ref={infiniteScrollElementRef}
            style={{
              height: '6px',
              width: '6px',
            }}
          />
        )}
        {loading && (
          <div className='flex flex-column mt-4'>
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
            <Skeleton height={'40px'} className='mb-3' />
          </div>
        )}

        {loggedActivities.map((e: any, index) => (
          <React.Fragment key={`${e.__typename}-${e.id}`}>
            {getTimelineItemByType(e.__typename, e, index)}
          </React.Fragment>
        ))}
        <div id={styles.scrollAnchor} />
      </div>
    </div>
  );
};
