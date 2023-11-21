'use client';

import React from 'react';
import { useParams } from 'next/navigation';

import { Box } from '@ui/layout/Box';
import { Contract } from '@graphql/types';
import { Select } from '@ui/form/SyncSelect';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { contractButtonSelect } from '@organization/src/components/Tabs/shared/contractSelectStyles';
import { ARRForecast } from '@organization/src/components/Tabs/panels/AccountPanel/ARRForecast/ARRForecast';

import { Notes } from './Notes';
import { EmptyContracts } from './EmptyContracts';
import { ContractCard } from './Contract/ContractCard';
import { AccountPanelSkeleton } from './AccountPanelSkeleton';
import { OrganizationPanel } from '../OrganizationPanel/OrganizationPanel';
import {
  useAccountPanelStateContext,
  AccountModalsContextProvider,
} from './context/AccountModalsContext';

export const AccountPanel = () => {
  const id = useParams()?.id as string;

  const { isModalOpen } = useAccountPanelStateContext();
  const client = getGraphQLClient();
  const { data, isInitialLoading, error } = useGetContractsQuery(client, {
    id,
  });
  if (isInitialLoading) {
    return <AccountPanelSkeleton />;
  }
  if (!data?.organization?.contracts?.length && error) {
    return <EmptyContracts name={data?.organization?.name || ''} />;
  }

  return (
    <AccountModalsContextProvider>
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
        {!!data?.organization?.contracts &&
          data?.organization?.contracts.map((contract) => (
            <>
              <ARRForecast
                opportunity={contract.opportunities?.[0]}
                name={data?.organization?.name || ''}
                isInitialLoading={isInitialLoading}
              />
              <ContractCard
                organizationId={id}
                organizationName={data?.organization?.name ?? ''}
                key={`contract-card-${contract.id}`}
                data={(contract as Contract) ?? undefined}
              />
            </>
          ))}

        <Notes id={id} data={data?.organization} />
      </OrganizationPanel>
    </AccountModalsContextProvider>
  );
};
