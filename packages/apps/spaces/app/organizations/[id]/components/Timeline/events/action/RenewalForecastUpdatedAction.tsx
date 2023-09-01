import React from 'react';
import { Flex } from '@ui/layout/Flex';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getFeatureIconColor } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { Text } from '@ui/typography/Text';
import { Action } from '@graphql/types';

const DEFAULT_COLOR_SCHEME = 'gray';

// You may want to tweak this function according to your needs
const getCurrencyString = (text: string) => {
  const match = text.split(/(\$[\d,]+(\.\d{2})?)/).filter(Boolean);
  return match?.[1];
};

interface RenewalForecastUpdatedActionProps {
  data: Action;
}

export const RenewalForecastUpdatedAction: React.FC<
  RenewalForecastUpdatedActionProps
> = ({ data }) => {
  let metadata;
  try {
    metadata = data?.metadata && JSON.parse(data.metadata);
  } catch (error) {
    console.error('Error parsing metadata', error);
  }

  const forecastedAmount = data.content && getCurrencyString(data.content);
  const [preText, postText] = data.content?.split('by ') ?? [];
  const isCreatedBySystem = data.content?.includes('default');

  const colorScheme =
    forecastedAmount && isCreatedBySystem
      ? getFeatureIconColor(metadata?.likelihood)
      : DEFAULT_COLOR_SCHEME;

  const authorText = isCreatedBySystem ? data.content : `${preText} by`;

  return (
    <Flex alignItems='center'>
      <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
        <Icons.Calculator />
      </FeaturedIcon>

      <Text
        my={1}
        maxW='500px'
        noOfLines={2}
        ml={2}
        fontSize='sm'
        color='gray.700'
      >
        {authorText}
        {!isCreatedBySystem && (
          <Text color='gray.500' as='span' ml={1}>
            {postText}
          </Text>
        )}
      </Text>
    </Flex>
  );
};
