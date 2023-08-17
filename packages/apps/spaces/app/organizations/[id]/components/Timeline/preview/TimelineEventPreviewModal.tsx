import React from 'react';
import { Card } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import styles from '../events/email/EmailPreviewModal.module.scss';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/TimelineEventsPreviewContext/TimelineEventPreviewContext';

import { EmailPreviewModal } from '@organization/components/Timeline/events/email/EmailPreviewModal';
import { SlackThreadPreviewModal } from '@organization/components/Timeline/events/slack/SlackThreadPreviewModal';

export const TimelineEventPreviewModal: React.FC<{
  invalidateQuery: () => void;
}> = ({ invalidateQuery }) => {
  const { closeModal, isModalOpen, modalContent } =
    useTimelineEventPreviewContext();

  if (!isModalOpen || !modalContent) {
    return null;
  }
  if (modalContent?.channel === 'EMAIL') {
    return <EmailPreviewModal invalidateQuery={invalidateQuery} />;
  }

  return (
    <div className={styles.container}>
      <div
        className={styles.backdrop}
        onClick={() => (isModalOpen ? closeModal() : null)}
      />
      <ScaleFade initialScale={0.9} in={isModalOpen} unmountOnExit>
        <Flex justifyContent='center'>
          <Card
            maxWidth={800}
            minWidth={544}
            w='full'
            zIndex={7}
            borderRadius='xl'
            height='100%'
            bg='gray.25'
            maxHeight='calc(100vh - 6rem)'
          >
            {modalContent?.channel === 'SLACK' && <SlackThreadPreviewModal />}
          </Card>
        </Flex>
      </ScaleFade>
    </div>
  );
};
