'use client';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Action, OnboardingStatus } from '@graphql/types';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { getMetadata } from '@organization/src/components/Timeline/events/action/utils';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { getColorScheme } from './util';

export const OnboardingStatusChangedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  const event = modalContent as Action;
  const metadata = getMetadata(event?.metadata);
  const reasonForChange = metadata?.comments;
  const status = metadata?.status as OnboardingStatus;
  const colorScheme = getColorScheme(status);

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name='Onboarding status changed'
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
            <Flag04 />
          </FeaturedIcon>
          <Text
            maxW='500px'
            noOfLines={2}
            ml={2}
            fontSize='sm'
            color='gray.700'
          >
            {event?.content}
          </Text>
        </CardBody>

        {metadata?.comments && (
          <CardFooter
            as={Flex}
            p='0'
            pt='3'
            mt='4'
            align='center'
            borderTop='1px solid'
            borderTopColor='gray.200'
          >
            <File02 color='gray.400' />
            <Text
              maxW='500px'
              noOfLines={2}
              ml={2}
              fontSize='sm'
              color='gray.500'
            >
              {reasonForChange}
            </Text>
          </CardFooter>
        )}
      </Card>
    </>
  );
};
