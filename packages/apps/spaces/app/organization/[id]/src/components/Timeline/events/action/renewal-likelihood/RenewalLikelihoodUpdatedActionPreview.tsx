import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { Action, RenewalLikelihoodProbability } from '@graphql/types';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { TimelineEventPreviewHeader } from '@organization/src/components/Timeline/preview/header/TimelineEventPreviewHeader';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { getMetadata, getLikelihoodDisplayData } from '../utils';

export const RenewalLikelihoodUpdatedActionPreview = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const event = modalContent as Action;

  const { preText, likelihood, author } = getLikelihoodDisplayData(
    event?.content ?? '',
  );
  const metadata = getMetadata(event?.metadata);
  const likelihoodFormatted =
    likelihood?.toUpperCase() as RenewalLikelihoodProbability;

  return (
    <>
      <TimelineEventPreviewHeader
        date={event?.createdAt}
        name={'Renewal likelihood'}
        onClose={closeModal}
        copyLabel='Copy link to this event'
      />
      <Card m={6} mt={3} p='4' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={getFeatureIconColor(likelihoodFormatted)}
          >
            <Icons.HeartActivity />
          </FeaturedIcon>
          <Text
            maxW='500px'
            noOfLines={2}
            ml={2}
            fontSize='sm'
            color='gray.700'
          >
            {preText}
            <Text as='span' fontWeight='semibold'>
              {likelihood}
            </Text>
            <Text color='gray.500' as='span' ml={1}>
              by {author}
            </Text>
          </Text>
        </CardBody>

        {likelihood && author && (
          <CardFooter p='0' as={Flex} flexDir='column'>
            <Divider my='4' />
            <Flex align='flex-start'>
              {metadata.reason ? (
                <Icons.File2 color='gray.400' />
              ) : (
                <Icons.FileCross viewBox='0 0 16 16' color='gray.400' />
              )}

              <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                {metadata.reason || 'No reason provided'}
              </Text>
            </Flex>
          </CardFooter>
        )}
      </Card>
    </>
  );
};
