import React from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

// todo use generated type after gql schema for service item is merged
const ServiceItem = ({
  data,
}: {
  data: { price: number; billed: string; quantity: number };
}) => {
  return (
    <Box
      _hover={{ '& button': { opacity: 1 } }}
      sx={{ '& button': { opacity: 0 } }}
    >
      <Flex justifyContent='space-between'>
        <Text>
          {data.quantity === 1 ? (
            `${formatCurrency(data.price ?? 0)} one-time`
          ) : (
            <>
              {data.quantity} {data.quantity > 1 ? 'licenses' : 'license'}
              <Text as='span' fontSize='sm' mx={1}>
                x
              </Text>
              {formatCurrency(data.price ?? 0)}/{data.billed}
            </>
          )}
        </Text>
        <IconButton
          transition='opacity 0.2s linear'
          size='xs'
          variant='ghost'
          aria-label='Remove service'
          color='gray.400'
          icon={<Delete boxSize='4' />}
        />
      </Flex>
      {/*<Text fontSize='sm' color='gray.500'>*/}
      {/*  Professional tier*/}
      {/*</Text>*/}
    </Box>
  );
};

// todo use generated type after gql schema for service item is merged
interface ServicesListProps {
  name?: string;
  data?: Array<{ id: string; price: number; billed: string; quantity: number }>; // todo when BE contract is available
}

export const ServicesList = ({ data }: ServicesListProps) => {
  return (
    <Flex flexDir='column' gap={1}>
      {data?.map((service, i) => (
        <React.Fragment key={`service-item-${service.id}`}>
          <ServiceItem data={service} />
          {data.length - 1 !== i && (
            <Divider w='full' orientation='horizontal' />
          )}
        </React.Fragment>
      ))}
    </Flex>
  );
};
