import { useRouter, useSearchParams } from 'next/navigation';
import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { TimelineEvent } from '../../types';
import { useDeepLinkToOpenModal } from './useDeeplinkToOpenModal';
import { useTimelineEventCachedData } from './useTimelineEventCachedData';

export const noop = () => undefined;

interface TimelineEventPreviewContextMethods {
  closeModal: () => void;
  openModal: (id: string) => void;
}
interface TimelineEventPreviewState {
  isModalOpen: boolean;
  modalContent: TimelineEvent | null;
}

const TimelineEventPreviewContext =
  createContext<TimelineEventPreviewContextMethods>({
    openModal: noop,
    closeModal: noop,
  });

const TimelineEventPreviewStateContext =
  createContext<TimelineEventPreviewState>({
    isModalOpen: false,
    modalContent: null,
  });

export const useTimelineEventPreviewMethodsContext = () => {
  return useContext(TimelineEventPreviewContext);
};
export const useTimelineEventPreviewStateContext = () => {
  return useContext(TimelineEventPreviewStateContext);
};

export const TimelineEventPreviewContextContextProvider = ({
  children,
  id = '',
}: PropsWithChildren<{
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
  const { handleFindTimelineEventInCache } = useTimelineEventCachedData();

  const handleDeleteParams = () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.delete('events');
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });
    router.push(`?${params}`);
  };

  useDeepLinkToOpenModal({
    modalContent,
    setModalContent,
    setIsModalOpen,
    handleDeleteParams,
  });
  const updateUrlAndPosition = (timelineEventId: string, id: string) => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('events', timelineEventId);
    router.push(`?${params}`);

    setLastActivePosition({ ...lastActivePosition, [id]: params.toString() });
  };

  const handleOpenModal = (timelineEventId: string) => {
    setIsModalOpen(true);

    const event = handleFindTimelineEventInCache(
      timelineEventId,
    ) as TimelineEvent;
    if (event) {
      setModalContent(event);
      updateUrlAndPosition(timelineEventId, id);
    }
  };
  const handleCloseModal = () => {
    if (!isModalOpen) return;
    setIsModalOpen(false);
    setModalContent(null);
    handleDeleteParams();
  };

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
    <TimelineEventPreviewContext.Provider
      value={{
        openModal: handleOpenModal,
        closeModal: handleCloseModal,
      }}
    >
      <TimelineEventPreviewStateContext.Provider
        value={{
          isModalOpen,
          modalContent,
        }}
      >
        {children}
      </TimelineEventPreviewStateContext.Provider>
    </TimelineEventPreviewContext.Provider>
  );
};
