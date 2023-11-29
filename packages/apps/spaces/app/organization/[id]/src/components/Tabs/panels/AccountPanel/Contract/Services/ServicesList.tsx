import React, { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { UpdateServiceModal } from './modals/UpdateServiceModal';
import { useUpdateServiceModalContext } from './../../context/AccountModalsContext';

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
      <Flex
        as='button'
        flexDir='column'
        cursor='pointer'
        onClick={() => onOpen(data)}
        _hover={{ '& button': { opacity: 1 } }}
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
        <Flex justifyContent='space-between'>
          <Text>
            {![BilledType.Usage, BilledType.Once].includes(data.billed) && (
              <>
                {data.quantity}
                <Text as='span' fontSize='sm' mx={1}>
                  ×
                </Text>
              </>
            )}

            {formatCurrency(data.price ?? 0)}
            {getBilledTypeLabel(data.billed)}
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
      </Flex>
    </>
  );
};

interface ServicesListProps {
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({ data }: ServicesListProps) => {
  const [isLocalOpen, setIsLocalOpen] = useState(false);
  const [selectedService, setSelectedService] = useState<
    ServiceLineItem | undefined
  >(undefined);
  const { modal } = useUpdateServiceModalContext();

  const handleOpenModal = (service: ServiceLineItem) => {
    setSelectedService(service);
    modal.onOpen();
    setIsLocalOpen(true);
  };

  return (
    <Flex flexDir='column' gap={1}>
      {data
        ?.filter(({ endedAt }) => !endedAt)
        ?.map((service, i) => (
          <React.Fragment key={`service-item-${service.id}`}>
            <ServiceItem data={service} onOpen={handleOpenModal} />
            {data?.length - 1 !== i && (
              <Divider w='full' orientation='horizontal' />
            )}
          </React.Fragment>
        ))}
      <UpdateServiceModal
        data={selectedService}
        isOpen={modal.isOpen && isLocalOpen}
        onClose={() => {
          modal.onClose();
          setIsLocalOpen(false);
        }}
      />
    </Flex>
  );
};
