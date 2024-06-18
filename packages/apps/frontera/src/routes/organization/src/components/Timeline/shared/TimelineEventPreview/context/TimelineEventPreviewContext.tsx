import { useSearchParams } from 'react-router-dom';
import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';

import { TimelineEvent } from '../../../types';
import { useDeepLinkToOpenModal } from './useDeeplinkToOpenModal';
import { useTimelineEventCachedData } from './useTimelineEventCachedData';

export const noop = () => undefined;

interface TimelineEventPreviewContextMethods {
  closeModal: () => void;
  openModal: (id: string) => void;
  handleOpenInvoice: (id: string) => void;
}
interface TimelineEventPreviewState {
  isModalOpen: boolean;
  modalContent: TimelineEvent | null;
}

const TimelineEventPreviewContext =
  createContext<TimelineEventPreviewContextMethods>({
    openModal: noop,
    closeModal: noop,
    handleOpenInvoice: noop,
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
  const store = useStore();

  const [lastActivePosition, setLastActivePosition] = useLocalStorage(
    `customeros-player-last-position`,
    { [id]: 'tab=about' },
  );
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<TimelineEvent | null>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const { handleFindTimelineEventInCache } = useTimelineEventCachedData();

  const handleDeleteParams = () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.delete('events');
    params.delete('invoice');
    setLastActivePosition({
      ...lastActivePosition,
      [id]: params.toString(),
    });

    setSearchParams(params);
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
    setSearchParams(params);

    setLastActivePosition({ ...lastActivePosition, [id]: params.toString() });
  };

  const handleOpenModal = (timelineEventId: string) => {
    setIsModalOpen(true);

    const cachedEvent = handleFindTimelineEventInCache(
      timelineEventId,
    ) as TimelineEvent;
    const storeEvent = store.timelineEvents
      .getByOrganizationId(id)
      ?.find((event) => event.id === timelineEventId)?.value;

    const event = store.demoMode ? storeEvent : cachedEvent;

    if (event) {
      setModalContent(event as TimelineEvent);
      updateUrlAndPosition(timelineEventId, id);
    }
  };

  // TODO refactor candidate added to open invoice in timeline preview modal
  const handleOpenInvoice = (timelineEventId: string) => {
    setIsModalOpen(true);

    setModalContent({ id: timelineEventId, __typename: 'Invoice' });
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('invoice', timelineEventId);
    setSearchParams(params);

    setLastActivePosition({ ...lastActivePosition, [id]: params.toString() });
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
        handleOpenInvoice, // todo remove me when invoice event are available as timeline events
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
