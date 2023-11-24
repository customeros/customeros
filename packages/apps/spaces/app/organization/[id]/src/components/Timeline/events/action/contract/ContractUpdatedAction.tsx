import React, { useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { iconsByStatus } from '@organization/src/components/Timeline/events/action/contract/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface ContractUpdatedActionProps {
  data: Action;
}

export const ContractUpdatedAction: React.FC<ContractUpdatedActionProps> = ({
  data,
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const state = useMemo(() => {
    return data?.content?.includes('live')
      ? 'live'
      : data?.content?.includes('renewed')
      ? 'renewed'
      : 'ended';
  }, [data?.content]);

  if (!data.content) return null;

  return (
    <Flex
      alignItems='center'
      onClick={() => openModal(data.id)}
      cursor='pointer'
    >
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
        {data.content}
        <Text as='span' fontWeight='semibold'>
          {state}
        </Text>
      </Text>
    </Flex>
  );
};
