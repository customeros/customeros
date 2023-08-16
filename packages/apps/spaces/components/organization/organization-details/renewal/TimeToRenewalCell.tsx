import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import formatDistanceToNow from 'date-fns/formatDistanceToNow';

export const TimeToRenewalCell = () => {
  return (
    <Text fontSize='sm' color='gray.700'>
      7 weeks
    </Text>
  );
};
