'use client';
import { useRouter, useSearchParams } from 'next/navigation';

import { useDebounce } from 'rooks';

import { Input } from '@ui/form/Input';
import { Flex } from '@ui/layout/Flex';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { UserPresence } from '@shared/components/UserPresence';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InputGroup, InputLeftElement } from '@ui/form/InputGroup';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

export const Search = () => {
  const client = getGraphQLClient();
  const router = useRouter();
  const searchParams = useSearchParams();
  const defaultValue = searchParams?.get('search') ?? '';
  const preset = searchParams?.get('preset');
  const { data: tenantNameQuery, isPending } = useTenantNameQuery(client);

  const placeholder =
    preset === 'customer'
      ? 'Search customers'
      : preset === 'portfolio'
      ? 'Search portfolio'
      : 'Search organizations';

  const handleChange = useDebounce(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const value = event.target.value;
      const params = new URLSearchParams(searchParams?.toString());

      if (!value) {
        params.delete('search');
      } else {
        params.set('search', value);
      }

      router.replace(`?${params}`);
    },
    300,
  );

  return (
    <Flex align='center' justify='space-between' pr='4'>
      <InputGroup w='full' size='lg' bg='gray.25' onChange={handleChange}>
        <InputLeftElement w='9'>
          <SearchSm boxSize='5' />
        </InputLeftElement>
        <Input
          pl='9'
          autoCorrect='off'
          spellCheck={false}
          placeholder={placeholder}
          defaultValue={defaultValue}
          borderBottom='unset'
          _hover={{
            borderBottom: 'unset',
          }}
          _focusWithin={{
            borderBottom: 'unset',
          }}
          _focus={{
            borderBottom: 'unset',
          }}
          _focusVisible={{
            borderBottom: 'unset',
          }}
        />
      </InputGroup>
      {!isPending && (
        <UserPresence channelName={`finder:${tenantNameQuery?.tenant}`} />
      )}
    </Flex>
  );
};
