import React, { useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody } from '@ui/presentation/Card';
import { iconsByStatus } from '@organization/src/components/Timeline/events/action/contract/utils';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

export const ContractUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;
  const state = useMemo(() => {
    return event?.content?.includes('live')
      ? 'live'
      : event?.content?.includes('renewed')
      ? 'renewed'
      : 'ended';
  }, [event?.content]);
  const contractName = event?.metadata;

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name={`${contractName} ${iconsByStatus[state].text} ${state}`}
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={iconsByStatus[state].colorScheme}
          >
            {iconsByStatus[state].icon}
          </FeaturedIcon>
          <Text
            my={1}
            maxW='500px'
            noOfLines={2}
            ml={2}
            fontSize='sm'
            color='gray.700'
          >
            {event?.content}
            <Text as='span' fontWeight='semibold'>
              {state}
            </Text>
          </Text>
        </CardBody>
      </Card>
    </>
  );
};
