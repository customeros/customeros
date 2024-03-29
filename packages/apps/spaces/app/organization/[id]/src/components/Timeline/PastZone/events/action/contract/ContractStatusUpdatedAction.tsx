import React, { useMemo } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Action } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { getMetadata } from '@organization/src/components/Timeline/PastZone/events/action/utils';
import { iconsByStatus } from '@organization/src/components/Timeline/PastZone/events/action/contract/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface ContractStatusUpdatedActionProps {
  data: Action;
}

export const ContractStatusUpdatedAction: React.FC<
  ContractStatusUpdatedActionProps
> = ({ data }) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const status = useMemo(() => {
    return getMetadata(data?.metadata)?.status?.toLowerCase();
  }, [data?.metadata]);

  // handle this situation
  if (!data.content || !status) return null;
  const content = data.content.replace(status.split('_').join(' '), '');

  return (
    <Flex
      alignItems='center'
      onClick={() => openModal(data.id)}
      cursor='pointer'
    >
      <FeaturedIcon
        size='md'
        minW='10'
        colorScheme={iconsByStatus[status]?.colorScheme as string}
      >
        {iconsByStatus[status]?.icon}
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
        <Text
          as='span'
          fontWeight={status === 'renewed' ? 'normal' : 'semibold'}
          ml={1}
        >
          {status.split('_').join(' ')}
        </Text>
      </Text>
    </Flex>
  );
};
