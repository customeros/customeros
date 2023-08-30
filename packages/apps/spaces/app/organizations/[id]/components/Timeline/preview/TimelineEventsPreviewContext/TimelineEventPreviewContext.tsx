import {
  useRef,
  useState,
  useEffect,
  RefObject,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';
import { useLocalStorage } from 'usehooks-ts';
import { useRouter, useSearchParams } from 'next/navigation';

import { InteractionEvent, Meeting } from '@graphql/types';

export const noop = () => undefined;

type Event = InteractionEvent | Meeting;

interface TimelineEventPreviewContextContextMethods {
  containerRef: RefObject<HTMLDivElement> | null;
  openModal: (content: Event) => void;
  closeModal: () => void;
  modalContent: Event | null;
  isModalOpen: boolean;
  events: Event[];
}

const TimelineEventPreviewContextContext =
  createContext<TimelineEventPreviewContextContextMethods>({
    containerRef: null,
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
}: PropsWithChildren<{ data: Event[]; id: string }>) => {
  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [id]: 'tab=about' },
  );

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<Event | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const searchParams = useSearchParams();

  const handleOpenModal = (content: Event) => {
    setIsModalOpen(true);
    const params = new URLSearchParams(searchParams ?? '');
    params.set('events', content.id);
    router.push(`?${params}`);
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });
    setModalContent(content);
  };

  const handleCloseModal = () => {
    if (!isModalOpen) return;
    const params = new URLSearchParams(searchParams ?? '');
    params.delete('events');
    setIsModalOpen(false);
    setModalContent(null);
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });
    router.push(`?${params}`);
  };

  useEffect(() => {
    const eventId = searchParams?.get('events');
    if (eventId && !modalContent) {
      const selectedEvent = data.find((d) => d.id === eventId);
      if (!selectedEvent) {
        // load more
        return;
      }
      setModalContent(selectedEvent);
      setIsModalOpen(true);
      // TODO: Load timeline event by ID and open modal
    }
  }, [searchParams]);

  return (
    <TimelineEventPreviewContextContext.Provider
      value={{
        containerRef,
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
