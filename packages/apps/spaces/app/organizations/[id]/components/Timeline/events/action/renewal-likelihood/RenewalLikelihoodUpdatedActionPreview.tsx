import React from 'react';
import { Card, CardFooter } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { TimelineEventPreviewHeader } from '@organization/components/Timeline/preview/header/TimelineEventPreviewHeader';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/context/TimelineEventPreviewContext';
import { CardBody } from '@chakra-ui/card';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getFeatureIconColor } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { Divider } from '@ui/presentation/Divider';
import { Action, Maybe, RenewalLikelihoodProbability } from '@graphql/types';
import { getLikelihoodDisplayData, getMetadata } from '../utils';

export const RenewalLikelihoodUpdatedActionPreview = () => {
  const { closeModal, modalContent } = useTimelineEventPreviewContext();
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
function getRenewalColor(
  data?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (data) {
    case 'HIGH':
      return 'success.500';
    case 'MEDIUM':
      return 'warning.500';
    case 'LOW':
      return 'error.500';
    case 'ZERO':
      return 'gray.700';
    default:
      return 'gray.400';
  }
}
