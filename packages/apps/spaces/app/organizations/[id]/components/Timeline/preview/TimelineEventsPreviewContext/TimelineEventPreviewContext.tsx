import { PropsWithChildren, RefObject, useContext } from 'react';
import { createContext, useState, useEffect, useRef } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { InteractionEvent } from '@graphql/types';

export const noop = () => undefined;

interface TimelineEventPreviewContextContextMethods {
  TimelineEventPreviewContextContainerRef: RefObject<HTMLDivElement> | null;
  openModal: (content: InteractionEvent) => void;
  closeModal: () => void;
  modalContent: InteractionEvent | null;
  isModalOpen: boolean;
  events: InteractionEvent[];
}

const TimelineEventPreviewContextContext =
  createContext<TimelineEventPreviewContextContextMethods>({
    TimelineEventPreviewContextContainerRef: null,
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
}: PropsWithChildren<{ data: InteractionEvent[] }>) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<InteractionEvent | null>(null);
  const TimelineEventPreviewContextContainerRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const handleOpenModal = (content: InteractionEvent) => {
    setIsModalOpen(true);

    setModalContent(content);
    const url = `${pathname}?events=${content.id}`;
    // Set URL parameter to ID of timeline event
    router.push(url);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setModalContent(null);
    // Clear URL parameters
    if (pathname) {
      router.replace(pathname);
    }
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
        TimelineEventPreviewContextContainerRef,
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
