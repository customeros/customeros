import { useSearchParams } from 'react-router-dom';

import debounce from 'lodash/debounce';

import { Input } from '@ui/form/Input/Input';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

export const Search = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const defaultValue = searchParams?.get('search') ?? '';

  const placeholder = 'Search organizations';

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    const params = new URLSearchParams(searchParams?.toString());

    if (!value) {
      params.delete('search');
    } else {
      params.set('search', value);
    }

    setSearchParams(params);
  };

  return (
    <div className='flex w-full items-center justify-between pr-4'>
      <InputGroup
        onChange={debounce(handleChange, 300)}
        className='w-full bg-gray-25 hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-2'
      >
        <LeftElement className='ml-2'>
          <SearchSm className='size-5' />
        </LeftElement>
        <Input
          size='lg'
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          placeholder={placeholder}
          defaultValue={defaultValue}
        />
      </InputGroup>
    </div>
  );
};
