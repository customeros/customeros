import { TimelineItem } from '@spaces/atoms/timeline-item';
import { NoteTimelineItem } from '@spaces/molecules/note-timeline-item';
import { EmailTimelineItemTemp } from '@spaces/molecules/conversation-timeline-item/EmailTimelineItemTemp';
import { PhoneConversationTimelineItem } from '@spaces/molecules/conversation-timeline-item/PhoneConversationTimelineItem';
import { ConversationTimelineItem } from '@spaces/molecules/conversation-timeline-item';
import { WebActionTimelineItem } from '@spaces/molecules/web-action-timeline-item';
import { InteractionTimelineItem } from '@spaces/molecules/interaction-timeline-item';
import { IssueTimelineItem } from '@spaces/molecules/issue-timeline-item';
import { EmailTimelineItem } from '@spaces/molecules/email-timeline-item';
import { LiveEventTimelineItem } from '@spaces/molecules/live-event-timeline-item';
import { MeetingTimelineItem } from '@spaces/molecules/meeting-timeline-item';
import React from 'react';

export const TimelineItemByType = ({
  type,
  data,
  index,
  loggedActivities,
  mode,
  contactName,
  id,
}: any) => {
  const getItem = () => {
    switch (type) {
      case 'Note':
        return (
          <TimelineItem
            source={data.source || data.appSource}
            first={index == 0}
            createdAt={data?.createdAt}
            externalLinks={data?.mentioned?.[0]?.externalLinks}
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
              participants={[...(data?.contacts || []), ...(data?.users || [])]}
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
          .filter((e: any) => e.__typename === 'Analysis')
          .filter((e: any) => e.analysisType !== 'summary')
          .find((e: any) => e.describes[0].id === data.describes[0].id);

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
            externalLinks={data.externalLinks}
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
              externalLinks={data?.issue?.externalLinks}
            >
              <EmailTimelineItem
                {...data}
                contactId={mode === 'CONTACT' && id}
              />
            </TimelineItem>
          );
        }
        if (data.channel === 'VOICE') {
          //we are using this to render the phone calls manually created by the user
          return (
            <ConversationTimelineItem
              id={data.id}
              content={undefined}
              transcript={[
                {
                  text: data.content,
                  party: {
                    from: data.sentBy,
                    to: data.sentTo,
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
  return <> {getItem()}</>;
};
