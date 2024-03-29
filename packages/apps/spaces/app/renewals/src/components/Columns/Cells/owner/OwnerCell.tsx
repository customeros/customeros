import { useEffect } from 'react';

import { produce } from 'immer';

import { User } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useRenewalsMeta } from '@shared/state/RenewalsMeta.atom';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName' | 'name'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
}

export const OwnerCell = ({ owner }: OwnerProps) => {
  const client = getGraphQLClient();
  const [renewalsMeta, setRenewalsMeta] = useRenewalsMeta();

  const { getUsers } = renewalsMeta;

  const { data } = useGetUsersQuery(
    client,
    {
      pagination: {
        limit: 1000,
        page: 1,
      },
    },
    {
      enabled: !getUsers.hasFetched,
    },
  );

  const value = data?.users?.content?.find((e) => e.id === owner?.id);
  const name =
    value?.name ??
    [owner?.firstName, owner?.lastName].filter(Boolean).join(' ').trim();

  useEffect(() => {
    if (!getUsers.hasFetched) {
      setRenewalsMeta(
        produce(renewalsMeta, (draft) => {
          draft.getUsers.hasFetched = false;
        }),
      );
    }
  }, [getUsers.hasFetched]);

  return (
    <Flex
      w='full'
      gap='1'
      align='center'
      _hover={{
        '& #edit-button': {
          opacity: 1,
        },
      }}
    >
      <Text cursor='default' color={name ? 'gray.700' : 'gray.400'}>
        {name || 'Owner'}
      </Text>
    </Flex>
  );
};
