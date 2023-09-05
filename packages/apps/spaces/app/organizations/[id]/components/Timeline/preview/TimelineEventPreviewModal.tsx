import { useState, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Card } from '@ui/presentation/Card';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { Action, InteractionEvent, Meeting } from '@graphql/types';

import { EmailPreviewModal } from '../events/email/EmailPreviewModal';
import { MeetingPreviewModal } from '../events/meeting/MeetingPreviewModal';
import { SlackThreadPreviewModal } from '../events/slack/SlackThreadPreviewModal';
import { ActionPreviewModal } from '../events/action/ActionPreviewModal';

import { useTimelineEventPreviewContext } from './TimelineEventsPreviewContext/TimelineEventPreviewContext';

interface TimelineEventPreviewModalProps {
  invalidateQuery: () => void;
}

export const TimelineEventPreviewModal = ({
  invalidateQuery,
}: TimelineEventPreviewModalProps) => {
  const [isMounted, setIsMounted] = useState(false); // needed for delaying the backdrop filter
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();

  useEffect(() => {
    setIsMounted(isModalOpen);
  }, [isModalOpen]);

  if (!isModalOpen || !modalContent) {
    return null;
  }

  const event = modalContent as InteractionEvent | Meeting | Action;
  const isMeeting = event?.__typename === 'Meeting';
  const isAction = event?.__typename === 'Action';
  const isInteraction = event?.__typename === 'InteractionEvent';
  const isSlack = isInteraction && event?.channel === 'SLACK';
  const isEmail = isInteraction && event?.channel === 'EMAIL';
  const handleCloseModal = () => {
    if (isEmail) return; // email modal handles closing the modal by itself
    closeModal();
  };

  return (
    <Flex
      position='absolute'
      top='0'
      bottom='0'
      left='0'
      right='0'
      zIndex={1}
      cursor='pointer'
      backdropFilter='blur(3px)'
      justify='center'
      background={isMounted ? 'rgba(16, 24, 40, 0.45)' : 'rgba(16, 24, 40, 0)'}
      align='center'
      transition='all 0.1s linear'
      onClick={handleCloseModal}
    >
      <ScaleFade
        in={isModalOpen}
        style={{
          position: 'absolute',
          marginInline: 'auto',
          top: '1rem',
          width: '544px',
          minWidth: '544px',
        }}
      >
        <Card
          size='lg'
          position='absolute'
          mx='auto'
          top='4'
          w='544px'
          minW='544px'
          cursor='default'
          onClick={(e) => e.stopPropagation()}
        >
          {isMeeting && (
            <MeetingPreviewModal invalidateQuery={invalidateQuery} />
          )}
          {isSlack && <SlackThreadPreviewModal />}
          {isEmail && <EmailPreviewModal invalidateQuery={invalidateQuery} />}
          {isAction && <ActionPreviewModal type={event.actionType} />}
        </Card>
      </ScaleFade>
    </Flex>
  );
};
