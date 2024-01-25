import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

function getBilledTypeLabel(billedType: BilledType): string {
  switch (billedType) {
    case BilledType.Annually:
      return '/year';
    case BilledType.Monthly:
      return '/month';
    case BilledType.None:
      return '';
    case BilledType.Once:
      return ' one-time';
    case BilledType.Usage:
      return '/use';
    case BilledType.Quarterly:
      return '/quarter';
    default:
      return '';
  }
}

const ServiceItem = ({
  data,
  onOpen,
}: {
  data: ServiceLineItem;
  onOpen: (props: ServiceLineItem) => void;
}) => {
  const allowedFractionDigits = data.billed === BilledType.Usage ? 4 : 2;

  return (
    <>
      <Flex
        w='full'
        as='button'
        flexDir='column'
        cursor='pointer'
        onClick={() => onOpen(data)}
        _focusVisible={{
          '&': {
            boxShadow: 'var(--chakra-shadows-outline)',
            outline: 'none',
            borderRadius: 'md',
          },
        }}
        sx={{ '& button': { opacity: 0 } }}
      >
        {data.name && (
          <Text fontSize='sm' color='gray.500' noOfLines={1} textAlign='left'>
            {data.name}
          </Text>
        )}
        <Flex justifyContent='space-between' w='full'>
          <Text>
            {![BilledType.Usage, BilledType.Once].includes(data.billed) && (
              <>
                {data.quantity}
                <Text as='span' fontSize='sm' mx={1}>
                  ×
                </Text>
              </>
            )}

            {formatCurrency(data.price ?? 0, allowedFractionDigits)}
            {getBilledTypeLabel(data.billed)}
          </Text>
        </Flex>
      </Flex>
    </>
  );
};

interface ServicesListProps {
  onModalOpen: () => void;
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({ data, onModalOpen }: ServicesListProps) => {
  const filteredData = data?.filter(({ endedAt }) => !endedAt);

  return (
    <Flex flexDir='column' gap={1}>
      {filteredData?.map((service, i) => (
        <React.Fragment key={`service-item-${service.id}`}>
          <ServiceItem data={service} onOpen={onModalOpen} />
          {filteredData?.length > 1 && filteredData?.length - 1 !== i && (
            <Divider w='full' orientation='horizontal' />
          )}
        </React.Fragment>
      ))}
    </Flex>
  );
};
