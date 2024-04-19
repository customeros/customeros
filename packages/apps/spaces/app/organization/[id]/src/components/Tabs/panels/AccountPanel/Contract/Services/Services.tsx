import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { ServiceLineItem } from '@graphql/types';
import { ServicesList } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/ServicesList';

interface Props {
  onModalOpen: () => void;
  currency?: string | null;
  data?: Array<ServiceLineItem> | null;
}

export const Services: React.FC<Props> = ({ data, currency, onModalOpen }) => {
  return (
    <>
      <Flex w='full' alignItems='center' justifyContent='space-between'>
        <Text fontWeight='semibold' fontSize='sm'>
          {!data?.length ? 'No services' : 'Services'}
        </Text>

        <IconButton
          size='xs'
          variant='ghost'
          aria-label={!data?.length ? 'Add services' : 'Edit services'}
          color='gray.400'
          onClick={() => {
            onModalOpen();
          }}
          icon={!data?.length ? <Plus /> : <Edit03 />}
        />
      </Flex>

      {data?.length && (
        <ServicesList
          data={data}
          onModalOpen={onModalOpen}
          currency={currency}
        />
      )}
    </>
  );
};
