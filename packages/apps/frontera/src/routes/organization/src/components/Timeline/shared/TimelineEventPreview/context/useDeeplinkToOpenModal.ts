import { useSearchParams } from 'react-router-dom';
import { Dispatch, useEffect, SetStateAction } from 'react';

import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { OrganizationQuery } from '@organization/graphql/organization.generated';
import { useTimelineEventCachedData } from '@organization/components/Timeline/shared/TimelineEventPreview/context/useTimelineEventCachedData';
import {
  GetTimelineEventsQuery,
  GetTimelineEventsDocument,
  useGetTimelineEventsQuery,
} from '@organization/graphql/getTimelineEvents.generated';

import { TimelineEvent } from '../../../types';

export const useDeepLinkToOpenModal = ({
  modalContent,
  setModalContent,
  setIsModalOpen,
  handleDeleteParams,
}: {
  handleDeleteParams: () => void;
  modalContent: TimelineEvent | null;
  setIsModalOpen: Dispatch<SetStateAction<boolean>>;
  setModalContent: Dispatch<SetStateAction<TimelineEvent | null>>;
}) => {
  const abortController = new AbortController();

  const [searchParams] = useSearchParams();
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
    const invoiceId = searchParams?.get('invoice');

    if (eventId && !modalContent && !invoiceId) {
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

    if (invoiceId) {
      setModalContent({
        id: invoiceId,
        __typename: 'Invoice',
      } as TimelineEvent);
      setIsModalOpen(true);
    }

    return () => {
      abortController.abort();
    };
  }, [searchParams]);
};
