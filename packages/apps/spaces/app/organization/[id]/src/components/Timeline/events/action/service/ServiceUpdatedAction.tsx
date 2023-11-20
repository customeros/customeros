import React, { useMemo } from 'react';

import { Dot } from '@ui/media/Dot';
import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { XCircle } from '@ui/media/icons/XCircle';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface ServiceUpdatedActionProps {
  data: Action;
}

export const ServiceUpdatedAction: React.FC<ServiceUpdatedActionProps> = ({
  data,
}) => {
  const colorScheme = useMemo(() => {
    return data?.content?.includes('added')
      ? 'primary'
      : data?.content?.includes('removed')
      ? 'error'
      : 'gray';
  }, [data?.content]);
  const { openModal } = useTimelineEventPreviewMethodsContext();
  if (!data.content) return null;

  return (
    <Flex
      alignItems='center'
      onClick={() => openModal(data.id)}
      cursor='pointer'
    >
      <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
        {data.content?.includes('removed') ? <XCircle /> : <Dot />}
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
