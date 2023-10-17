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

export const useDeepLinkToOpenModal = (
  modalContent: TimelineEvent | null,
  setModalContent: Dispatch<SetStateAction<TimelineEvent | null>>,
  setIsModalOpen: Dispatch<SetStateAction<boolean>>,
  handleDeleteParams: () => void,
) => {
  const searchParams = useSearchParams();
  const client = getGraphQLClient();
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

      queryClient.setQueryData<OrganizationQuery>(
        singleEventQueryKey,
        (oldData) => {
          return result;
        },
      );

      if (!result.timelineEvents.length) {
        handleDeleteParams();
        toastError(
          "Sorry, we couldn't find this event",
          `timeline-event-not-found-${id}`,
        );
      }
      return result.timelineEvents[0] as TimelineEvent;
    } catch (error) {
      handleDeleteParams();
      toastError(
        "Sorry, we couldn't find this event",
        `timeline-event-not-found-${id}`,
      );
    }
  };

  useEffect(() => {
    const eventId = searchParams?.get('events');
    if (eventId && !modalContent) {
      // Assuming that handleFindEventInCache and getModalContentFromServer functions are available in this scope

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
      setModalContent(selectedEvent);
      setIsModalOpen(true);
    }
  }, [searchParams]);
};
