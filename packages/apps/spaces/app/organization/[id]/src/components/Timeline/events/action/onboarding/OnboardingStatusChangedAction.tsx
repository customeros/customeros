import React, { useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Action, OnboardingStatus } from '@graphql/types';
import { getMetadata } from '@organization/src/components/Timeline/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import { getColorScheme } from './util';

const statusMap = {
  [OnboardingStatus.Late]: 'Late',
  [OnboardingStatus.OnTrack]: 'On track',
  [OnboardingStatus.Done]: 'Done',
  [OnboardingStatus.Stuck]: 'Stuck',
  [OnboardingStatus.NotStarted]: 'Not started',
  [OnboardingStatus.NotApplicable]: 'Not applicable',
  [OnboardingStatus.Successful]: 'Successful',
};

interface OnboardingStatusChangedActionProps {
  data: Action;
}

export const OnboardingStatusChangedAction: React.FC<
  OnboardingStatusChangedActionProps
> = ({ data }) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const status = useMemo(() => {
    return getMetadata(data?.metadata)?.status;
  }, [data?.metadata]) as OnboardingStatus;

  // handle this situation
  if (!data.content || !status) return null;

  const statusLabel = statusMap[status];
  const content = data?.content.replace(statusLabel, '').trimEnd();
  const colorScheme = getColorScheme(status);

  return (
    <Flex
      alignItems='center'
      onClick={() => openModal(data.id)}
      cursor='pointer'
    >
      <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
        <Flag04 />
      </FeaturedIcon>

      <Text
        my={1}
        maxW='500px'
        noOfLines={2}
        ml={2}
        fontSize='sm'
        color='gray.700'
      >
        {content}
        <Text as='span' fontWeight='semibold' ml={1}>
          {statusLabel}
        </Text>
      </Text>
    </Flex>
  );
};
