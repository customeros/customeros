import { useRef, useMemo, useEffect } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { EditSkuUsecase } from '@domain/usecases/settings-products/edit-sku.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { MaskedInput } from '@ui/form/Input/MaskedInput.tsx';
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
    editSkuUsecase.resetErrors();
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
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  editSkuUsecase.errors.productName,
              })}
              onChange={(e) => {
                editSkuUsecase.editProductName(e.target.value);

                if (e.target.value && editSkuUsecase.errors.productName) {
                  editSkuUsecase.resetNameError();
                }
              }}
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
              <p className='text-xs text-error-600 pl-2.5'>
                {editSkuUsecase.errors.productName}
              </p>
            )}
          </div>

          <div className='flex flex-col'>
            <label htmlFor={'sku-price'} className='text-sm font-medium mb-1'>
              Price
            </label>
            <MaskedInput
              size='sm'
              mask={'num'}
              id='sku-price'
              variant='outline'
              data-test='sku-price'
              placeholder='Price per unit'
              onFocus={(e) => e.target.select()}
              value={editSkuUsecase.maskedPrice}
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
              onAccept={(v, instance) => {
                editSkuUsecase.editPrice(instance?.unmaskedValue || '');
                editSkuUsecase.setMaskedValue(v);

                if (v.length && editSkuUsecase.errors.price) {
                  editSkuUsecase.resetPriceError();
                }
              }}
              blocks={{
                num: {
                  mask: Number,
                  scale: 2,
                  lazy: true,
                  min: 0,
                  radix: '.',
                  placeholderChar: '#',
                  thousandsSeparator: ',',
                  normalizeZeros: true,
                  padFractionalZeros: true,
                  autofix: false,
                },
              }}
            />
            {editSkuUsecase.errors.price && (
              <p className='text-xs text-error-600 pl-2.5'>
                {editSkuUsecase.errors.price}
              </p>
            )}

            {!editSkuUsecase.errors.price && (
              <p className='text-xs text-grayModern-500 pl-2.5'>
                Product currency is set in your company's contracts
              </p>
            )}
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
            data-test='edit-sku-confirm'
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
