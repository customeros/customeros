import { useRef, useMemo, useEffect } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { AddSkuUsecase } from '@domain/usecases/settings-products/add-sku.usecase';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { SkuType } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Command,
  CommandCancelButton,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const AddNewSku = observer(() => {
  const { ui } = useStore();
  const inputRef = useRef<HTMLInputElement>(null);
  const addSkuUsecase = useMemo(() => {
    return new AddSkuUsecase();
  }, []);

  const handleConfirm = async () => {
    addSkuUsecase.resetError();
    addSkuUsecase.validate();

    if (addSkuUsecase.errors.price || addSkuUsecase.errors.productName) {
      return;
    }

    addSkuUsecase.createSku();
  };

  useModKey('Enter', () => {
    ui.commandMenu.setOpen(false);
  });

  useEffect(() => {
    const focusTimer = setTimeout(() => {
      if (inputRef.current) {
        inputRef.current.focus({ preventScroll: true });
      }
    }, 0);

    return () => clearTimeout(focusTimer);
  }, []); // Run only on mount

  const handleClose = () => {
    addSkuUsecase.reset();
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
          <h1 className='text-base font-semibold inline'>New product</h1>
          <CommandCancelIconButton onClose={handleClose} />
        </div>

        <div className={'mt-4 gap-2.5 flex flex-col'}>
          <div className='flex flex-col'>
            <label htmlFor='sku-type' className='text-sm font-medium mb-1.5'>
              Type
            </label>

            <RadioGroup
              id={'sku-type'}
              name='sku-type'
              className={'gap-2'}
              value={addSkuUsecase.type}
              onValueChange={(val: SkuType) => addSkuUsecase.editType(val)}
            >
              <Radio value={SkuType.Subscription}>
                <span className='text-sm'>Subscription</span>
              </Radio>
              <Radio value={SkuType.OneTime}>
                <span className='text-sm'>One-time</span>
              </Radio>
            </RadioGroup>
          </div>
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
              value={addSkuUsecase.productName}
              placeholder='Name of product offering'
              onChange={(e) => {
                addSkuUsecase.editProductName(e.target.value);
              }}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  addSkuUsecase.errors.productName,
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

            {addSkuUsecase.errors.productName && (
              <p className='text-xs text-error-600'>
                {addSkuUsecase.errors.productName}
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
              value={addSkuUsecase.price}
              placeholder='Price per unit'
              onChange={(e) => {
                addSkuUsecase.editPrice(e.target.value);
              }}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  addSkuUsecase.errors.price,
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
            {addSkuUsecase.errors.price && (
              <p className='text-xs text-error-600'>
                {addSkuUsecase.errors.price}
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
            // isLoading={addSkuUsecase.isValidating}
            data-test='contact-actions-confirm-flow-change'
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleConfirm();
              }
            }}
          >
            Create product
          </Button>
        </div>
      </article>
    </Command>
  );
});
