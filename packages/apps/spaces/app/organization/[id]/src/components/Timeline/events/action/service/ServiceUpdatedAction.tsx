import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface ServiceUpdatedActionProps {
  data: Action;
  mode?: 'created' | 'updated';
}

export const ServiceUpdatedAction: React.FC<ServiceUpdatedActionProps> = ({
  data,
  mode = 'updated',
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  if (!data.content) return null;
  const isTemporary = data.appSource === 'customeros-optimistic-update';

  return (
    <Flex
      alignItems='center'
      opacity={isTemporary ? 0.5 : 1}
      onClick={() => !isTemporary && openModal(data.id)}
      cursor={isTemporary ? 'progress' : 'pointer'}
    >
      <FeaturedIcon
        size='md'
        minW='10'
        colorScheme={mode === 'created' ? 'primary' : 'gray'}
      >
        <DotSingle />
      </FeaturedIcon>

      <Text
        my={1}
        maxW='500px'
        noOfLines={2}
        ml={2}
        fontSize='sm'
        color='gray.700'
      >
        {data.content}
      </Text>
    </Flex>
  );
};
