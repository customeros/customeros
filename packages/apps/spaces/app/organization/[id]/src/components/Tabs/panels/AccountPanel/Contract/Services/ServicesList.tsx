import { useParams } from 'next/navigation';
import React, { useRef, useState } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { toastError } from '@ui/presentation/Toast';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useCloseServiceLineItemMutation } from '@organization/src/graphql/closeServiceLineItem.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';

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
    case BilledType.Quarterly:
      return '/quarter';
    default:
      return '';
  }
}

const ServiceItem = ({
  data,
  onOpen,
  onCloseService,
}: {
  data: ServiceLineItem;
  onOpen: (props: ServiceLineItem) => void;
  onCloseService: (props: { input: { id: string } }) => void;
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
          <IconButton
            transition='opacity 0.2s linear'
            size='xs'
            variant='ghost'
            aria-label='Remove service'
            color='gray.400'
            icon={<Delete boxSize='4' />}
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              onCloseService({ input: { id: data?.id } });
            }}
          />
        </Flex>
      </Flex>
    </>
  );
};

interface ServicesListProps {
  contractId: string;
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({ data, contractId }: ServicesListProps) => {
  const [isLocalOpen, setIsLocalOpen] = useState(false);
  const [selectedService, setSelectedService] = useState<
    ServiceLineItem | undefined
  >(undefined);
  const { modal } = useUpdateServiceModalContext();
  const orgId = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryKey = useGetContractsQuery.getKey({ id: orgId });
  const closeServiceLineItem = useCloseServiceLineItemMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.id === contractId,
          );
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return {
                ...contractData,
                serviceLineItems: contractData?.serviceLineItems?.filter(
                  (el) => el.id !== input.id,
                ),
              };
            });
          }
        });
      });

      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );
      toastError('Failed to remove service', 'remove-service-error');
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
      }, 1000);
    },
  });

  const handleOpenModal = (service: ServiceLineItem) => {
    setSelectedService(service);
    modal.onOpen();
    setIsLocalOpen(true);
  };

  const filteredData = data?.filter(({ endedAt }) => !endedAt);

  return (
    <Flex flexDir='column' gap={1}>
      {filteredData?.map((service, i) => (
        <React.Fragment key={`service-item-${service.id}`}>
          <ServiceItem
            data={service}
            onOpen={handleOpenModal}
            onCloseService={closeServiceLineItem.mutate}
          />
          {filteredData?.length > 1 && filteredData?.length - 1 !== i && (
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
