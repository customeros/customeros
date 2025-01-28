import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

import { ProductRow } from './ProductRow';

export const ProductsList = observer(
  ({ searchTerm }: { searchTerm: string }) => {
    const store = useStore();
    const skusArray = store.skus.toArray();

    const data =
      searchTerm.trim().length > 0
        ? new Fuse(skusArray, {
            keys: [{ name: 'name', getFn: (o) => o.value.name }],
            threshold: 0.3,
            isCaseSensitive: false,
          })
            .search(searchTerm)
            .map((r) => r.item)
        : skusArray;

    if (!data.length) {
      return (
        <div className='text-sm text-center text-grayModern-500 '>
          No products in sight...
        </div>
      );
    }

    return (
      <div>
        <div className='grid grid-cols-[minmax(0,1fr)_minmax(0,108px)_minmax(0,118px)_28px] w-full text-sm'>
          <div className='font-medium flex items-center'>Product</div>
          <div className='font-medium flex items-center'>Type</div>
          <div className='font-medium flex items-center'>Price</div>
          <div className='w-7 h-[28px]' />
        </div>

        {data.map((row) => (
          <ProductRow id={row.id} key={`product-${row.value.id}`} />
        ))}
      </div>
    );
  },
);
