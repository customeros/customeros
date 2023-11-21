import React, { useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useUpdateServiceModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { UpdateServiceModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/UpdateServiceModal';

function getBilledTypeLabel(billedType: BilledType): string {
  switch (billedType) {
    case BilledType.Annually:
      return '/year';
    case BilledType.Monthly:
      return '/month';
    case BilledType.None:
      return '';
    case BilledType.Once:
      return '/one-time';
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
  return (
    <>
      <Box
        cursor='pointer'
        onClick={() => onOpen(data)}
        _hover={{ '& button': { opacity: 1 } }}
        sx={{ '& button': { opacity: 0 } }}
      >
        {data.name && (
          <Text fontSize='sm' color='gray.500'>
            {data.name}
          </Text>
        )}
        <Flex justifyContent='space-between'>
          <Text>
            {data.billed === BilledType.Once ? (
              `${formatCurrency(data.price ?? 0)} one-time`
            ) : (
              <>
                {data.quantity} {data.quantity > 1 ? 'licenses' : 'license'}
                <Text as='span' fontSize='sm' mx={1}>
                  x
                </Text>
                {formatCurrency(data.price ?? 0)}
                {getBilledTypeLabel(data.billed)}
              </>
            )}
          </Text>
          {/*<IconButton*/}
          {/*  transition='opacity 0.2s linear'*/}
          {/*  size='xs'*/}
          {/*  variant='ghost'*/}
          {/*  aria-label='Remove service'*/}
          {/*  color='gray.400'*/}
          {/*  icon={<Delete boxSize='4' />}*/}
          {/*/>*/}
        </Flex>
      </Box>
    </>
  );
};

// todo use generated type after gql schema for Services item is merged
interface ServicesListProps {
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({ data }: ServicesListProps) => {
  const [selectedService, setSelectedService] = useState<
    ServiceLineItem | undefined
  >(undefined);
  const { modal } = useUpdateServiceModalContext();

  const handleOpenModal = (service: ServiceLineItem) => {
    setSelectedService(service);
    modal.onOpen();
  };

  return (
    <Flex flexDir='column' gap={1}>
      {data?.map((service, i) => (
        <React.Fragment key={`service-item-${service.id}`}>
          <ServiceItem data={service} onOpen={handleOpenModal} />
          {data.length - 1 !== i && (
            <Divider w='full' orientation='horizontal' />
          )}
        </React.Fragment>
      ))}
      <UpdateServiceModal
        data={selectedService}
        isOpen={modal.isOpen}
        onClose={modal.onClose}
      />
    </Flex>
  );
};
