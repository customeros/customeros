import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { SearchSm } from '@ui/media/icons/SearchSm.tsx';

import { ProductsList } from './components/ProductsList';

export const Products = observer(() => {
  const store = useStore();
  const [searchTerm, setSearchTerm] = useState('');

  const hasProducts = store.skus.value.size > 0;

  return (
    <div className='pt-2 max-w-[500px] border-r border-gray-200 h-full'>
      <div className='flex flex-col gap-4 px-6 pb-4 '>
        <div>
          <div className='flex justify-between'>
            <p
              data-test='products-header'
              className='text-gray-700  font-semibold'
            >
              Products
            </p>
            <Button
              size='xs'
              leftIcon={<Plus />}
              colorScheme='primary'
              dataTest={'add-product-button'}
              onClick={() => {
                store.ui.commandMenu.setType('AddNewSku');
                store.ui.commandMenu.setOpen(true);
              }}
            >
              Product
            </Button>
          </div>

          <p className='text-sm'>
            Create and manage product offerings used for invoicing
          </p>
        </div>

        <div className='flex items-center gap-2'>
          <SearchSm />
          <Input
            size='sm'
            className='w-full'
            placeholder='Search for products'
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        {!hasProducts && (
          <p className='text-sm text-center text-grayModern-500 '>
            Nothing to search for yet. Go ahead, add your first product...
          </p>
        )}
      </div>
      {hasProducts && <ProductsList searchTerm={searchTerm} />}
    </div>
  );
});
