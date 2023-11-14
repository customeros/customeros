import React from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';

// todo fix type when BE contract is available
const ServiceItem = ({ data }: { data: unknown }) => {
  return (
    <Box>
      <Flex justifyContent='space-between'>
        <Text>10 licenses @ $200/month, per license</Text>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Remove service'
          color='gray.400'
          icon={<Delete boxSize='4' />}
        />
      </Flex>
      <Text fontSize='sm' color='gray.500'>
        Professional tier
      </Text>
    </Box>
  );
};

interface ServicesListProps {
  name?: string;
  data?: Array<unknown>; // todo when BE contract is available
}

export const ServicesList = ({ data }: ServicesListProps) => {
  return (
    <Flex flexDir='column'>
      {data?.map((d) => (
        // @ts-expect-error TODO: fix type when BE contract is available
        <ServiceItem data={d} key={`service-item-${d.id}`} />
      ))}
    </Flex>
  );
};
