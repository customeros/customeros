import { Flex } from '@ui/layout/Flex';
import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';
import { Ticket02 } from '@ui/media/icons/Ticket02';
import { FC } from 'react';
export const EmptyIssueMessage: FC<{ organizationName: string }> = ({
  organizationName,
}) => (
  <Flex direction='column' alignItems='center' mt='4'>
    <Box
      border='1px solid'
      borderColor='gray.200'
      padding={3}
      borderRadius='md'
      mb={6}
    >
      <Ticket02 color='gray.700' boxSize='6' />
    </Box>
    <Text color='gray.700' fontWeight={600}>
      No issues detected
    </Text>
    <Text color='gray.500' mt={1} mb={6} textAlign='center'>
      It looks like {organizationName} has had a smooth journey thus far. Or
      perhaps they’ve been shy about reporting issues. Stay proactive and keep
      monitoring for optimal support.
    </Text>
  </Flex>
);
