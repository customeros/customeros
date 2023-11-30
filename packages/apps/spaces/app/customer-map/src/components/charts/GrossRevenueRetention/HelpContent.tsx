'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        Gross Revenue Retention (GRR) tells you what percentage of revenue you
        keep from all your customers over their lifetime.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        To determine this, we compare the Monthly Recurring Revenue (MRR) from
        the current period (minus any up-sells and cross-sells) with the initial
        contracted MRR.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        This comparison is then expressed as a percentage, with a maximum value
        of 100% indicating that all original revenue has been retained.
      </Text>
    </Box>
  );
};
