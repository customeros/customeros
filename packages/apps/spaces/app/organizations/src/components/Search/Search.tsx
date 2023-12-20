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
      size='lg'
      bg='gray.25'
      onChange={debounce(handleChange, 300)}
    >
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
  );
};
