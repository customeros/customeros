import { useRouter, useSearchParams } from 'next/navigation';

import debounce from 'lodash/debounce';

import { Input } from '@ui/form/Input';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { InputGroup, InputLeftElement } from '@ui/form/InputGroup';

export const Search = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const defaultValue = searchParams?.get('search') ?? '';
  const preset = searchParams?.get('preset');

  const placeholder =
    preset === 'customer'
      ? 'Search customers'
      : preset === 'portfolio'
      ? 'Search portfolio'
      : 'Search organizations';

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    const params = new URLSearchParams(searchParams?.toString());

    if (!value) {
      params.delete('search');
    } else {
      params.set('search', value);
    }

    router.push(`?${params}`);
  };

  return (
    <InputGroup
      w='full'
      mb='2'
      size='lg'
      bg='gray.25'
      borderRadius='1rem'
      border='1px solid'
      borderColor='gray.200'
      onChange={debounce(handleChange, 300)}
    >
      <InputLeftElement>
        <SearchSm boxSize='6' />
      </InputLeftElement>
      <Input
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
  );
};
