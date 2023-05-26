import React, { useEffect, useRef, useState } from 'react';

import { ConversationTimelineItem } from '@spaces/molecules/conversation-timeline-item';
import { EmailTimelineItem } from '@spaces/molecules/email-timeline-item';
import { LiveEventTimelineItem } from '@spaces/molecules/live-event-timeline-item';
import { NoteTimelineItem } from '@spaces/molecules/note-timeline-item';
import { WebActionTimelineItem } from '@spaces/molecules/web-action-timeline-item';
import { IssueTimelineItem } from '@spaces/molecules/issue-timeline-item';
import { EmailTimelineItemTemp } from '@spaces/molecules/conversation-timeline-item/EmailTimelineItemTemp';
import { PhoneConversationTimelineItem } from '@spaces/molecules/conversation-timeline-item/PhoneConversationTimelineItem';
import { MeetingTimelineItem } from '@spaces/molecules//meeting-timeline-item';
import { InteractionTimelineItem } from '@spaces/molecules/interaction-timeline-item';
import {
  NoActivityTimelineElement,
  TimelineItem,
  TimelineItemSkeleton,
} from '@spaces/atoms/timeline-item';
import { useInfiniteScroll } from './useInfiniteScroll';
import classNames from 'classnames';

import styles from './timeline.module.scss';
import { AnimatePresence } from 'framer-motion';

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
  const timelineContainerRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef(null);
  const [useAnchoring, setUseAnchoring] = useState(true);

  const anchor = useRef<HTMLDivElement>(null);

  const infiniteScrollElementRef = useRef(null);
  useInfiniteScroll({
    element: infiniteScrollElementRef,
    isFetching: loading,
    callback: () => {
      if (loggedActivities.length > 10) {
        onLoadMore(timelineContainerRef);
      }
    },
  });

  useEffect(() => {
    if (anchor.current && useAnchoring && !loading && loggedActivities.length) {
      anchor?.current?.scrollIntoView();
      setTimeout(() => {
        setUseAnchoring(false);
      }, 1000);
    }
  }, [anchor, loggedActivities, useAnchoring, loading]);

  const getTimelineItemByType = (type: string, data: any, index: number) => {
    switch (type) {
      case 'Note':
        return (
          <TimelineItem
            source={data.source || data.appSource}
            first={index == 0}
            createdAt={data?.createdAt}
          >
            <NoteTimelineItem note={data} />
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
        if (data.describes.find((e: any) => e.__typename === 'Meeting')) {
          return null;
        }
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
            source={data.source}
            content={decodeContent(data.content)}
            transcript={decodeContent(transcriptForSummary.content)}
            type={data.analysisType}
            createdAt={data?.createdAt}
            contentType={transcriptForSummary.contentType}
            mode='PHONE_CALL' // fixme - mode will be assessed from data inside the component (on message base)
          />
        );
      }
      case 'PageView':
        return (
          <TimelineItem
            source={data.source}
            first={index == 0}
            createdAt={data?.startedAt}
          >
            <WebActionTimelineItem {...data} contactName={contactName} />
          </TimelineItem>
        );
      case 'InteractionSession':
        return (
          <TimelineItem
            source={data.source}
            first={index == 0}
            createdAt={data?.startedAt}
          >
            <InteractionTimelineItem
              {...data}
              contactId={contactName && id}
              organizationId={!contactName && id}
            />
          </TimelineItem>
        );
      case 'Issue':
        return (
          <TimelineItem
            source={data.source}
            first={index == 0}
            createdAt={data?.createdAt}
          >
            <IssueTimelineItem {...data} />
          </TimelineItem>
        );

      case 'InteractionEvent':
        if (data.channel === 'EMAIL') {
          return (
            <TimelineItem
              source={data.source}
              first={index == 0}
              createdAt={data?.createdAt}
            >
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
              source={data.source}
              type={'summary'} //fixme: this is used to get the same style as the summary of  phone call
              createdAt={data?.createdAt}
              mode='PHONE_CALL' // fixme - mode will be assessed from data inside the component (on message base)
            />
          );
        } else {
          return null;
        }
      case 'LiveEventTimelineItem':
        return (
          <LiveEventTimelineItem
            first={index == 0}
            contactId={id}
            source={data.source}
          />
        );
      case 'Meeting':
        return (
          <TimelineItem
            source={data.source}
            first={index == 0}
            createdAt={data?.createdAt || new Date()}
            hideTimeTick
          >
            <MeetingTimelineItem meeting={data} />
          </TimelineItem>
        );
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
      <div
        className={classNames(styles.timelineContent, styles.scrollable)}
        ref={containerRef}
      >
        <div
          ref={infiniteScrollElementRef}
          style={{
            height: '1px',
            width: '100%',
            display: useAnchoring ? 'none' : 'block',
          }}
        />
        <AnimatePresence mode='wait'>
          {loading && (
            <>
              <TimelineItemSkeleton key='timeline-element-skeleton-1' />
              <TimelineItemSkeleton key='timeline-element-skeleton-2' />
            </>
          )}
          {noActivity && (
            <NoActivityTimelineElement key='no-activity-timeline-item' />
          )}

          {loggedActivities.map((e: any, index) => {
            return (
              <React.Fragment
                key={`${e.__typename}-${e.id}-${index}-timeline-element`}
              >
                {getTimelineItemByType(e.__typename, e, index)}
              </React.Fragment>
            );
          })}
          <LiveEventTimelineItem
            key='live-stream-timeline-item'
            first={false}
            contactId={id}
            source={'LiveStream'}
          />
          <div
            className={styles.scrollAnchor}
            ref={anchor}
            key='chat-scroll-timeline-anchor'
          />
        </AnimatePresence>
      </div>
    </div>
  );
};
