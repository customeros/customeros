import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { ArchiveSkuUsecase } from '@domain/usecases/settings-products/archive-sku.usecase';

import { cn } from '@ui/utils/cn';
import { Edit03 } from '@ui/media/icons/Edit03';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

const archiveUsecase = new ArchiveSkuUsecase();
export const ProductRow = observer(({ id }: { id: string }) => {
  const {
    onOpen: onOpenArchive,
    onClose,
    open: openDelete,
  } = useDisclosure({ id: 'delete-field' });
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false);
  const store = useStore();

  const row = store.skus.getById(id);

  if (!row || !row?.id) return null;

  return (
    <div className='grid grid-cols-[minmax(40px,1fr)_minmax(0,100px)_minmax(0,110px)_40px] w-full text-sm group gap-x-2'>
      <div title={row.value.name} className='truncate flex items-center'>
        <span className='truncate'>{row.value.name}</span>
      </div>
      <div className='truncate flex items-center'>{row.typeLabel}</div>
      <div className='truncate flex items-center justify-end'>
        {row.formattedPrice}
      </div>
      <div className='w-7'>
        <Menu onOpenChange={(state) => setIsMenuOpen(state)}>
          <MenuButton asChild>
            <IconButton
              size='xs'
              variant='ghost'
              icon={<DotsVertical />}
              aria-label='Edit product'
              className={cn('opacity-0 group-hover:opacity-100', {
                'opacity-100': isMenuOpen,
              })}
            />
          </MenuButton>
          <MenuList>
            <MenuItem
              className='group/edit'
              onClick={() => {
                store.ui.commandMenu.setType('EditSku');
                store.ui.commandMenu.setContext({
                  ...store.ui.commandMenu.context,
                  ids: [row?.id as string],
                });
                store.ui.commandMenu.setOpen(true);
              }}
            >
              <div className='flex items-center'>
                <Edit03 className='mr-2 text-gray-500 group-hover/edit:text-gray-700' />
                Edit product
              </div>
            </MenuItem>
            <MenuItem onClick={onOpenArchive} className='group/archive'>
              <div className='flex items-center'>
                <Archive className='mr-2 group-hover/archive:text-gray-700 text-gray-500' />
                Archive product
              </div>
            </MenuItem>
          </MenuList>
        </Menu>
      </div>
      <ConfirmDeleteDialog
        onClose={onClose}
        isOpen={openDelete}
        label={`Archive ${row.value.name}?`}
        confirmButtonLabel='Archive product'
        description={`Archiving this product will NOT remove it from existing invoices. It will only discontinue it and stop others from using it for future invoices.`}
        onConfirm={() => {
          if (row?.id) {
            archiveUsecase.archiveSku(row.id);

            return;
          }
          console.error('No SKU id found');
        }}
      />
    </div>
  );
});
