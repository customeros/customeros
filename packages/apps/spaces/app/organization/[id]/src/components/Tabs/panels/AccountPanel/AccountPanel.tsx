'use client';

import React, { FC, PropsWithChildren } from 'react';
import { useParams, useRouter } from 'next/navigation';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Select } from '@ui/form/SyncSelect';
import { Organization } from '@graphql/types';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { contractButtonSelect } from '@organization/src/components/Tabs/shared/contractSelectStyles';
import { Contracts } from '@organization/src/components/Tabs/panels/AccountPanel/Contracts/Contracts';

import { Notes } from './Notes';
import { EmptyContracts } from './EmptyContracts';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import {
  useAccountPanelStateContext,
  AccountModalsContextProvider,
} from './context/AccountModalsContext';

const AccountPanelComponent = () => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const router = useRouter();

  const { isModalOpen } = useAccountPanelStateContext();
  const { data, isFetching } = useGetContractsQuery(client, {
    id,
  });

  if (isFetching) {
    return <AccountPanelSkeleton />;
  }

  if (!data?.organization?.contracts?.length) {
    return (
      <EmptyContracts name={data?.organization?.name || ''}>
        <Notes id={id} data={data?.organization} />
      </EmptyContracts>
    );
  }

  return (
    <OrganizationPanel
      title='Account'
      withFade
      actionItem={
        <Box display='none'>
          <Select
            isSearchable={false}
            isClearable={false}
            isMulti={false}
            value={{
              label: 'Customer',
              value: 'customer',
            }}
            options={[
              {
                label: 'Customer',
                value: 'customer',
              },
              {
                label: 'Prospect',
                value: 'prospect',
              },
            ]}
            chakraStyles={{
              ...contractButtonSelect,
              container: (props, state) => {
                const isCustomer = state.getValue()[0]?.value === 'customer';

                return {
                  ...props,
                  px: 2,
                  pointerEvents: 'none',
                  py: '1px',
                  border: '1px solid',
                  borderColor: isCustomer ? 'success.200' : 'gray.300',
                  backgroundColor: isCustomer ? 'success.50' : 'transparent',
                  color: isCustomer ? 'success.700' : 'gray.500',

                  borderRadius: '2xl',
                  fontSize: 'xs',
                  maxHeight: '22px',

                  '& > div': {
                    p: 0,
                    border: 'none',
                    fontSize: 'xs',
                    maxHeight: '22px',
                    minH: 'auto',
                  },
                };
              },
              valueContainer: (props, state) => {
                const isCustomer = state.getValue()[0]?.value === 'customer';

                return {
                  ...props,
                  p: 0,
                  border: 'none',
                  fontSize: 'xs',
                  maxHeight: '22px',
                  minH: 'auto',
                  color: isCustomer ? 'success.700' : 'gray.500',
                };
              },
              singleValue: (props) => {
                return {
                  ...props,
                  maxHeight: '22px',
                  p: 0,
                  minH: 'auto',
                  color: 'inherit',
                };
              },
              menuList: (props) => {
                return {
                  ...props,
                  w: 'fit-content',
                  left: '-32px',
                };
              },
            }}
            leftElement={<ActivityHeart color='success.500' mr='1' />}
          />
        </Box>
      }
      shouldBlockPanelScroll={isModalOpen}
    >
      <Contracts
        isFetching={isFetching}
        organization={data?.organization as Organization}
      />
      <Flex mb={8} />
      <Button
        position='fixed'
        left={'200px'}
        right={0}
        bottom={0}
        borderRadius={0}
        bg='gray.25'
        p={7}
        justifyContent='space-between'
        alignItems='center'
        rightIcon={<ChevronRight boxSize={4} color='gray.400' />}
        variant='ghost'
        _hover={{
          bg: 'gray.25',
          '& svg': {
            color: 'gray.500',
          },
        }}
        onClick={() => router.push(`?tab=invoices`)}
      >
        <Text fontSize='sm' fontWeight='semibold'>
          Invoices • 2
        </Text>
      </Button>
    </OrganizationPanel>
  );
};

export const AccountPanel: FC<PropsWithChildren> = () => (
  <AccountModalsContextProvider>
    <AccountPanelComponent />
  </AccountModalsContextProvider>
);
