'use client';

import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';

export const HelpContent = () => {
  return (
    <Box mt='1'>
      <Text fontSize='md' fontWeight='normal'>
        The customer map tells the story of all your customers.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        The circle sizes reflect the contracted ARR each of them brings to the
        table. The colour shows how healthy they are right now.
      </Text>
      <br />
      <Text fontSize='md' fontWeight='normal'>
        Their position on the timeline marks the signing date of their contract
        with you, laying out a history of when they came on board.
      </Text>
    </Box>
  );
};
