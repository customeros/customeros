import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface ServiceUpdatedActionProps {
  data: Action;
}

export const ServiceUpdatedAction: React.FC<ServiceUpdatedActionProps> = ({
  data,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  if (!data.content) return null;

  return (
    <Flex
      alignItems='center'
      onClick={() => openModal(data.id)}
      cursor='pointer'
    >
      <FeaturedIcon size='md' minW='10' colorScheme='gray'>
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
