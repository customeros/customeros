import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useLocalStorage } from 'usehooks-ts';
import { TimelineEvent } from '../../types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  GetTimelineEventDocument,
  GetTimelineEventQuery,
} from '@organization/src/graphql/getTimelineEvent.generated';
import { toastError } from '@ui/presentation/Toast';

export const noop = () => undefined;

interface TimelineEventPreviewContextContextMethods {
  openModal: (content: TimelineEvent) => void;
  closeModal: () => void;
  modalContent: TimelineEvent | null;
  isModalOpen: boolean;
  events: TimelineEvent[];
}

const TimelineEventPreviewContextContext =
  createContext<TimelineEventPreviewContextContextMethods>({
    openModal: noop,
    closeModal: noop,
    modalContent: null,
    isModalOpen: false,
    events: [],
  });

export const useTimelineEventPreviewContext = () => {
  return useContext(TimelineEventPreviewContextContext);
};

export const TimelineEventPreviewContextContextProvider = ({
  children,
  data = [],
  id = '',
}: PropsWithChildren<{
  data: TimelineEvent[];
  id: string;
}>) => {
  const client = getGraphQLClient();
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [id]: 'tab=about' },
  );

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<TimelineEvent | null>(null);
  const router = useRouter();
  const searchParams = useSearchParams();

  const handleOpenModal = (content: TimelineEvent) => {
    setIsModalOpen(true);
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('events', content.id);
    router.push(`?${params}`);
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });
    setModalContent(content);
  };

  const handleDeleteParams = () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.delete('events');
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });
    router.push(`?${params}`);
  };

  const handleCloseModal = () => {
    if (!isModalOpen) return;
    setIsModalOpen(false);
    setModalContent(null);
    handleDeleteParams();
  };

  const getModalContentFromServer = async (id: string) => {
    try {
      const result = await client.request<GetTimelineEventQuery>(
        GetTimelineEventDocument,
        {
          ids: [id],
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
      const selectedEvent = data.find((d) => d.id === eventId);
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
      // TODO: Load timeline event by ID and open modal
    }
  }, [searchParams]);

  useEffect(() => {
    const isEmail =
      modalContent?.__typename === 'InteractionEvent' &&
      modalContent?.channel === 'EMAIL';

    const handleCloseOnEsc = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handleCloseModal();
      }
    };
    if (isModalOpen && !isEmail) {
      document.addEventListener('keydown', handleCloseOnEsc);
    }
    if (!isModalOpen && !isEmail) {
      document.removeEventListener('keydown', handleCloseOnEsc);
    }
    return () => {
      document.removeEventListener('keydown', handleCloseOnEsc);
    };
  }, [isModalOpen]);

  return (
    <TimelineEventPreviewContextContext.Provider
      value={{
        openModal: handleOpenModal,
        closeModal: handleCloseModal,
        isModalOpen,
        modalContent,
        events: data,
      }}
    >
      {children}
    </TimelineEventPreviewContextContext.Provider>
  );
};
