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

  const handleCloseModal = () => {
    if (!isModalOpen) return;
    const params = new URLSearchParams(searchParams?.toString() ?? '');
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
