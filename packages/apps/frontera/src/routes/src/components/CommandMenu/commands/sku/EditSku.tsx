import { useRef, useMemo, useEffect } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { EditSkuUsecase } from '@domain/usecases/settings-products/edit-sku.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const EditSku = observer(() => {
  const { ui } = useStore();
  const inputRef = useRef<HTMLInputElement>(null);
  const contextId = ui.commandMenu.context.ids[0];
  const editSkuUsecase = useMemo(() => {
    return new EditSkuUsecase(contextId);
  }, []);

  const handleConfirm = async () => {
    editSkuUsecase.resetError();
    editSkuUsecase.validate();

    if (editSkuUsecase.errors.productName || editSkuUsecase.errors.price) {
      return;
    }
    editSkuUsecase.editSku();
    ui.commandMenu.setOpen(false);
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  useEffect(() => {
    if (editSkuUsecase) {
      editSkuUsecase.setInitial();
    }
  }, [editSkuUsecase]);

  useEffect(() => {
    const focusTimer = setTimeout(() => {
      if (inputRef.current) {
        inputRef.current.focus({ preventScroll: true });
      }
    }, 0);

    return () => clearTimeout(focusTimer);
  }, []); // Run only on mount

  const handleClose = () => {
    // editSkuUsecase.reset();
    ui.commandMenu.toggle('AddNewSku');
    ui.commandMenu.clearCallback();
  };

  useKeyBindings({
    Escape: handleClose,
  });

  return (
    <Command shouldFilter={false} className={'!w-auto'}>
      <article className={'p-6'}>
        <div className='flex justify-between items-center'>
          <h1 className='text-base font-semibold inline'>Edit product</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className={'mt-4 gap-2.5 flex flex-col'}>
          <p className='text-sm'>
            Editing a product’s details will not affect existing invoices, only
            future ones.
          </p>
          <div className='flex flex-col'>
            <label
              htmlFor={'sku-product-name'}
              className='text-sm font-medium mb-1'
            >
              Product name
            </label>
            <Input
              autoFocus
              size={'sm'}
              ref={inputRef}
              variant={'outline'}
              id='sku-product-name'
              dataTest='sku-product-name'
              value={editSkuUsecase.productName}
              placeholder='Name of product offering'
              onChange={(e) => {
                editSkuUsecase.editProductName(e.target.value);
              }}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  editSkuUsecase.errors.productName,
              })}
              onKeyDownCapture={(e) => {
                e.stopPropagation();

                if (e.key === 'Enter') {
                  handleConfirm();
                }

                if (e.key === 'Escape') {
                  handleClose();
                }
              }}
            />

            {editSkuUsecase.errors.productName && (
              <p className='text-xs text-error-600'>
                {editSkuUsecase.errors.productName}
              </p>
            )}
          </div>

          <div className='flex flex-col'>
            <label htmlFor={'sku-price'} className='text-sm font-medium mb-1'>
              Price
            </label>
            <Input
              size={'sm'}
              id='sku-price'
              type={'number'}
              variant={'outline'}
              dataTest='sku-price'
              value={editSkuUsecase.price}
              placeholder='Price per unit'
              onChange={(e) => {
                editSkuUsecase.editPrice(e.target.value);
              }}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  editSkuUsecase.errors.price,
              })}
              onKeyDownCapture={(e) => {
                e.stopPropagation();

                if (e.key === 'Enter') {
                  handleConfirm();
                }

                if (e.key === 'Escape') {
                  handleClose();
                }
              }}
            />
            {editSkuUsecase.errors.price && (
              <p className='text-xs text-error-600'>
                {editSkuUsecase.errors.price}
              </p>
            )}
            <p className='text-xs text-grayModern-500'>
              Product currency is set in your organization's contracts
            </p>
          </div>
        </div>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton onClose={handleClose} />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            onClick={handleConfirm}
            dataTest={'add-domain'}
            loadingText={'Adding domain...'}
            // isLoading={editSkuUsecase.isValidating}
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Update product
          </Button>
        </div>
      </article>
    </Command>
  );
});
