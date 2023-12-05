import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { Card, CardBody } from '@ui/presentation/Card';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

export const ServiceUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name='Service updated'
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' minW='10' colorScheme='gray'>
            <DotSingle />
          </FeaturedIcon>
          <Text
            maxW='500px'
            noOfLines={2}
            ml={2}
            fontSize='sm'
            color='gray.700'
          >
            {event.content}
          </Text>
        </CardBody>
      </Card>
    </>
  );
};
