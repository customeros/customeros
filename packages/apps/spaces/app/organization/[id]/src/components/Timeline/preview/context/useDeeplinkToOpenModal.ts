import * as Sentry from '@sentry/nextjs';
import { useSearchParams } from 'next/navigation';
import { Dispatch, SetStateAction, useEffect } from 'react';
import {
  GetTimelineEventsDocument,
  GetTimelineEventsQuery,
  useGetTimelineEventsQuery,
} from '@organization/src/graphql/getTimelineEvents.generated';
import { OrganizationQuery } from '@organization/src/graphql/organization.generated';
import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useQueryClient } from '@tanstack/react-query';
import { useTimelineEventCachedData } from '@organization/src/components/Timeline/preview/context/useTimelineEventCachedData';
import { TimelineEvent } from '../../types';

export const useDeepLinkToOpenModal = ({
  modalContent,
  setModalContent,
  setIsModalOpen,
  handleDeleteParams,
}: {
  modalContent: TimelineEvent | null;
  setModalContent: Dispatch<SetStateAction<TimelineEvent | null>>;
  setIsModalOpen: Dispatch<SetStateAction<boolean>>;
  handleDeleteParams: () => void;
}) => {
  const abortController = new AbortController();

  const searchParams = useSearchParams();
  const client = getGraphQLClient({ signal: abortController.signal });
  const queryClient = useQueryClient();
  const { handleFindTimelineEventInCache } = useTimelineEventCachedData();

  const getModalContentFromServer = async (id: string) => {
    try {
      const singleEventQueryKey = useGetTimelineEventsQuery.getKey({
        ids: [id],
      });
      const result = await client.request<GetTimelineEventsQuery>(
        GetTimelineEventsDocument,
        {
          ids: [id],
        },
      );

      queryClient.setQueryData<OrganizationQuery>(singleEventQueryKey, () => {
        return result;
      });

      if (!result.timelineEvents.length) {
        handleDeleteParams();
        toastError(
          "We couldn't find this event",
          `timeline-event-not-found-${id}`,
        );
      }
      return result.timelineEvents[0] as TimelineEvent;
    } catch (error) {
      Sentry.captureException(`Event not found: ${error}`);
      handleDeleteParams();
      toastError(
        "We couldn't find this event",
        `timeline-event-not-found-${id}`,
      );
      return null;
    }
  };

  useEffect(() => {
    const eventId = searchParams?.get('events');
    if (eventId && !modalContent) {
      const selectedEvent = handleFindTimelineEventInCache(eventId);

      if (!selectedEvent) {
        getModalContentFromServer(eventId).then((content) => {
          if (content) {
            setModalContent(content);
            setIsModalOpen(true);
          }
        });
        return;
      }
      setModalContent(selectedEvent as TimelineEvent);
      setIsModalOpen(true);
    }

    return () => {
      abortController.abort();
    };
  }, [searchParams]);
};
