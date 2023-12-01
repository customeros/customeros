import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { OpportunityRenewalLikelihood } from '@graphql/types';

import { getLikelihoodColor, getRenewalLikelihoodLabel } from './utils';

interface RenewalLikelihoodCellProps {
  value?: OpportunityRenewalLikelihood | null;
}

export const RenewalLikelihoodCell = ({
  value,
}: RenewalLikelihoodCellProps) => {
  return (
    <Flex key={Math.random()} w='full'>
      <Text
        cursor='default'
        color={value ? getLikelihoodColor(value) : 'gray.400'}
      >
        {value ? getRenewalLikelihoodLabel(value) : 'Unknown'}
      </Text>
    </Flex>
  );
};
