import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { ServiceLineItem } from '@graphql/types';
import { ServicesList } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/ServicesList';
import { useAddServiceModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { CreateServiceModal } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/CreateServiceModal';

interface Props {
  contractId: string;
  data?: Array<ServiceLineItem> | null;
}

export const Services: React.FC<Props> = ({ contractId, data }) => {
  const { modal } = useAddServiceModalContext();

  return (
    <>
      <Flex w='full' alignItems='center' justifyContent='space-between'>
        <Text fontWeight='semibold' fontSize='sm'>
          {!data?.length ? 'No services' : 'Services'}
        </Text>

        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Add service'
          color='gray.400'
          onClick={() => modal.onOpen()}
          icon={<Plus boxSize='4' />}
        />
      </Flex>

      {data?.length && <ServicesList data={data} />}
      <CreateServiceModal
        contractId={contractId}
        isOpen={modal.isOpen}
        onClose={modal.onClose}
      />
    </>
  );
};
