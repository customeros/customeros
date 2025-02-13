import Fuse from 'fuse.js';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

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
      <div className={'pl-6'}>
        <div className='grid grid-cols-[minmax(40px,1fr)_minmax(40px,100px)_minmax(0,110px)_40px] w-full text-sm gap-x-2'>
          <div className='font-medium flex items-center'>Product</div>
          <div className='font-medium flex items-center'>Type</div>
          <div className='font-medium flex items-center justify-end'>Price</div>
          <div className='w-8 h-[28px]' />
        </div>

        <ScrollAreaRoot>
          <ScrollAreaViewport style={{ height: 'calc(100vh - 150px)' }}>
            <div>
              {data.map((row) => (
                <ProductRow id={row.id} key={`product-${row.value.id}`} />
              ))}
            </div>
          </ScrollAreaViewport>
          <ScrollAreaScrollbar orientation='vertical'>
            <ScrollAreaThumb />
          </ScrollAreaScrollbar>
        </ScrollAreaRoot>
      </div>
    );
  },
);
