import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { XCircle } from '@ui/media/icons/XCircle';
import { Action, BilledType } from '@graphql/types';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getMetadata } from '@organization/src/components/Timeline/PastZone/events/action/utils';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
  const formattedContent = formatString(data.content, metadata?.billedType);

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

function formatString(str: string, type: string) {
  const digitCount = type === BilledType.Usage ? 4 : 2;
  const regex =
    type === BilledType.Usage ? /\b(\d+\.\d{4})\b/g : /\b(\d+\.\d{2})\b/g;

  return str.replace(regex, (_, number) => {
    return formatCurrency(Number(number), digitCount);
  });
}
