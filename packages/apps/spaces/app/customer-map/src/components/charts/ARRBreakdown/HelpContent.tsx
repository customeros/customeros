'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        Annual Recurring Revenue (ARR) is the total amount of recurring money
        you could expect to receive from your customers over the next 12 months.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        It includes income from new subscriptions, renewals, and any additional
        revenue from service upsells. However, it also accounts for reductions
        due to service downgrades, cancellations or lost customers and excludes
        one-time and per-use services.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        This gives you a clear picture of what drives your revenue.
      </Text>
    </Box>
  );
};
