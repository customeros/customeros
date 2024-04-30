import { useSearchParams } from 'react-router-dom';

import { useDebounce } from 'rooks';

import { Input } from '@ui/form/Input/Input';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { UserPresence } from '@shared/components/UserPresence';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';
import { useTenantNameQuery } from '@shared/graphql/tenantName.generated';

export const Search = () => {
  const client = getGraphQLClient();
  const [searchParams, setSeatchParams] = useSearchParams();
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

      setSeatchParams(params);
    },
    300,
  );

  return (
    <div className='flex items-center justify-between pr-4'>
      <InputGroup
        className='w-full bg-gray-25 hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-2'
        onChange={handleChange}
      >
        <LeftElement className='ml-2'>
          <SearchSm className='size-5' />
        </LeftElement>
        <Input
          size='lg'
          autoCorrect='off'
          spellCheck={false}
          placeholder={placeholder}
          defaultValue={defaultValue}
          variant='unstyled'
        />
      </InputGroup>
      {!isPending && (
        <UserPresence channelName={`finder:${tenantNameQuery?.tenant}`} />
      )}
    </div>
  );
};
