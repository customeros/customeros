import { useSearchParams } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';

import { CustomFieldModal } from '../CustomFieldModal';

interface HeaderProps {
  title: string;
  subTitle: string;
  numberOfCoreFields: number;
  numberOfCustomFields: number;
}

export const Header = ({
  title,
  subTitle,
  numberOfCoreFields = 0,
  numberOfCustomFields = 0,
}: HeaderProps) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { onOpen, onToggle, open } = useDisclosure();

  const handleItemClick = (tab: string) => () => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    params.set('view', tab);
    setSearchParams(params.toString());
  };
  const checkIsActive = (tab: string) => searchParams?.get('view') === tab;

  const checkIsActiveCustom = checkIsActive('custom');
  const checkIsActiveCore = checkIsActive('core');

  const dynamicClassesCustom = cn(
    checkIsActiveCustom
      ? ['font-medium', 'bg-white']
      : ['font-normal', 'bg-transparent', 'text-gray-500'],
  );

  const dynamicClassesCore = cn(
    checkIsActiveCore
      ? ['font-medium', 'bg-white']
      : ['font-normal', 'bg-transparent', 'text-gray-500'],
  );

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    params.set('search', e.target.value);
    setSearchParams(params.toString());
  };

  return (
    <>
      <div className='flex items-center justify-between pb-2 pt-[5px] sticky top-0 bg-gray-25 z-10'>
        <h1 className='font-medium'>{title}</h1>
        <Button
          size='xs'
          leftIcon={<Plus />}
          colorScheme='primary'
          onClick={() => onOpen()}
        >
          Custom field
        </Button>
      </div>
      <h2 className='text-sm'>{subTitle}</h2>
      <div className='flex flex-col gap-4 mt-4'>
        <ButtonGroup>
          <Button
            size='sm'
            onClick={handleItemClick('custom')}
            className={`w-[50%] ${dynamicClassesCustom} !border-r-[1px]`}
          >
            Custom • {numberOfCustomFields}
          </Button>
          <Button
            size='sm'
            onClick={handleItemClick('core')}
            className={`w-[50%] ${dynamicClassesCore}`}
          >
            Core • {numberOfCoreFields}
          </Button>
        </ButtonGroup>
        <InputGroup className=''>
          <LeftElement>
            <SearchSm className='text-gray-500' />
          </LeftElement>
          <Input
            size='sm'
            variant='unstyled'
            placeholder='Search fields...'
            onChange={(e) => handleSearch(e)}
            value={searchParams?.get('search') || ''}
          />
        </InputGroup>
      </div>
      {open && <CustomFieldModal isOpen={open} onOpenChange={onToggle} />}
    </>
  );
};
