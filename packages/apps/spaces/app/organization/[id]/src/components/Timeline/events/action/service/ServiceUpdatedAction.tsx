import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { XCircle } from '@ui/media/icons/XCircle';
import { Action, BilledType } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getMetadata } from '@organization/src/components/Timeline/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface ServiceUpdatedActionProps {
  data: Action;
  mode?: 'created' | 'updated' | 'removed';
}

export const ServiceUpdatedAction: React.FC<ServiceUpdatedActionProps> = ({
  data,
  mode = 'updated',
}) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const metadata = getMetadata(data?.metadata);
  if (!data.content) return null;
  const isTemporary = data.appSource === 'customeros-optimistic-update';
  const formattedContent = data.content
    .replace(
      metadata?.price,
      formatCurrency(
        Number(metadata?.price),
        metadata?.billedType === BilledType.Usage ? 4 : 2,
      ),
    )
    .replace(
      metadata?.previousPrice,
      formatCurrency(
        Number(metadata?.previousPrice),
        metadata?.billedType === BilledType.Usage ? 4 : 2,
      ),
    );

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
        {mode === 'removed' ? <XCircle /> : <DotSingle />}
      </FeaturedIcon>

      <Text
        my={1}
        maxW='500px'
        noOfLines={2}
        ml={2}
        fontSize='sm'
        color='gray.700'
      >
        {formattedContent}
      </Text>
    </Flex>
  );
};
