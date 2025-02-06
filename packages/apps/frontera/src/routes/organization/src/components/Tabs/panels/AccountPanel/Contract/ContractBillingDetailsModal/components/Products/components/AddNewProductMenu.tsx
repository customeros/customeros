import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus';
import { Combobox } from '@ui/form/Combobox';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { SkuType, BilledType } from '@graphql/types';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

interface AddNewProductMenuProps {
  contractId: string;
}
export const AddNewProductMenu = observer(
  ({ contractId }: AddNewProductMenuProps) => {
    const store = useStore();
    const contractLineItemsStore = store.contractLineItems;
    const [isOpen, setIsOpen] = useState(false);

    const options = store.skus.toArray().map((sku) => ({
      label: sku.value.name,
      value: sku.id,
      sku: sku,
    }));

    return (
      <Popover open={isOpen} onOpenChange={(open) => setIsOpen(open)}>
        <PopoverTrigger>
          <IconButton
            size='xxs'
            icon={<Plus />}
            className='ml-1'
            variant='outline'
            colorScheme='gray'
            aria-label='Add a product'
            dataTest='contract-card-add-sli'
          />
        </PopoverTrigger>
        <PopoverContent
          align='end'
          side='bottom'
          onClick={(e) => e.stopPropagation()}
          onMouseDown={(e) => e.stopPropagation()}
          className='py-1 min-w-[254px] max-w-[420px] z-[99999999]'
        >
          <Combobox
            escapeClearsValue
            options={options}
            closeMenuOnSelect
            placeholder={'Search for a product'}
            onKeyDown={(e) => {
              if (e.key === 'Escape') setIsOpen(false);
            }}
            noOptionsMessage={() => (
              <div className='py-2 cursor-default'>
                {options.length
                  ? 'No products found.'
                  : 'No products yet. Create them in Settings.'}
              </div>
            )}
            formatOptionLabel={(option) => (
              <div className='inline-flex items-center'>
                <span title={option.label} className='truncate max-w-[280px]'>
                  {option.label}
                </span>
                <span className='mx-0.5'>•</span>
                <span className='text-gray-500'>{option.sku.typeLabel}</span>
              </div>
            )}
            onChange={(newValue) => {
              contractLineItemsStore.create({
                billingCycle:
                  newValue?.sku.value.type === SkuType.Subscription
                    ? BilledType.Monthly
                    : BilledType.Once,
                contractId,
                skuId: newValue?.sku.id,
                price: newValue?.sku.value.price,
                serviceStarted: new Date(Date.now() + 86400000).toISOString(),
              });
              setIsOpen(false);
            }}
          />
        </PopoverContent>
      </Popover>
    );
  },
);
