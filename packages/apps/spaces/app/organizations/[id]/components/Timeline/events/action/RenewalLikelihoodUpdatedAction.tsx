import React from 'react';
import { Flex } from '@ui/layout/Flex';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { getFeatureIconColor } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { Text } from '@ui/typography/Text';
import { Action, RenewalLikelihoodProbability } from '@graphql/types';

const getLikelihoodDisplayData = (text: string) => {
  const match = text.match(/(.+? to )(.+?)(?: by )(.+)/);

  if (!match) {
    return { preText: '', likelihood: '', author: '' };
  }

  return {
    preText: match?.[1], // "Renewal likelihood set to "
    likelihood: match?.[2], // "Low"
    author: match?.[3], // " by Olivia Rhye"
  };
};

interface RenewalForecastUpdatedActionProps {
  data: Action;
}

export const RenewalLikelihoodUpdatedAction: React.FC<
  RenewalForecastUpdatedActionProps
> = ({ data }) => {
  if (!data.content) return null;
  const { preText, likelihood, author } = getLikelihoodDisplayData(
    data.content,
  );
  return (
    <Flex alignItems='center'>
      <FeaturedIcon
        size='md'
        minW='10'
        colorScheme={getFeatureIconColor(
          likelihood.toUpperCase() as RenewalLikelihoodProbability,
        )}
      >
        <Icons.HeartActivity />
      </FeaturedIcon>

      <Text
        my={1}
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
          {author}
        </Text>
      </Text>
    </Flex>
  );
};
