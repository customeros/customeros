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
import { MaskedInput } from '@ui/form/Input/MaskedInput.tsx';
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
    addSkuUsecase.resetErrors();
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
          <h1
            data-test='new-product-header'
            className='text-base font-semibold inline'
          >
            New product
          </h1>
          <CommandCancelIconButton dataTest='add-sku-x' onClose={handleClose} />
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
                <span
                  className='text-sm'
                  data-test={'add-product-subscription'}
                >
                  Subscription
                </span>
              </Radio>
              <Radio value={SkuType.OneTime}>
                <span className='text-sm' data-test={'add-product-one-time'}>
                  One-time
                </span>
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
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  addSkuUsecase.errors.productName,
              })}
              onChange={(e) => {
                addSkuUsecase.editProductName(e.target.value);

                if (e.target.value.length && addSkuUsecase.errors.price) {
                  addSkuUsecase.resetNameError();
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

            {addSkuUsecase.errors.productName && (
              <p className='text-xs text-error-600 pl-2.5'>
                {addSkuUsecase.errors.productName}
              </p>
            )}
          </div>

          <div className='flex flex-col'>
            <label htmlFor={'sku-price'} className='text-sm font-medium mb-1'>
              Price
            </label>

            <MaskedInput
              size='sm'
              mask={`num`}
              id='sku-price'
              variant='outline'
              dataTest='sku-price'
              placeholder='Price per unit'
              onFocus={(e) => (e.target as HTMLInputElement).select()}
              className={cn({
                'border-error-600 hover:!border-error-600 focus:!border-error-600 active:!border-error-600':
                  addSkuUsecase.errors.price,
              })}
              onAccept={(v, instance) => {
                addSkuUsecase.editPrice(instance?.unmaskedValue || '');

                if (v.length && addSkuUsecase.errors.price) {
                  addSkuUsecase.resetPriceError();
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
              blocks={{
                num: {
                  mask: Number,
                  scale: 2,
                  lazy: false,
                  min: 0,
                  radix: '.',
                  placeholderChar: '#',
                  thousandsSeparator: ',',
                  normalizeZeros: true,
                  padFractionalZeros: true,
                  autofix: true,
                },
              }}
            />

            {addSkuUsecase.errors.price && (
              <p className='text-xs text-error-600 pl-2.5'>
                {addSkuUsecase.errors.price}
              </p>
            )}

            {!addSkuUsecase.errors.price && (
              <p className='text-xs text-grayModern-500 pl-2.5'>
                Product currency is set in your company's contracts
              </p>
            )}
          </div>
        </div>

        <div className='flex justify-between gap-3 mt-6'>
          <CommandCancelButton
            onClose={handleClose}
            dataTest='add-sku-cancel'
          />

          <Button
            size='sm'
            variant='outline'
            className='w-full'
            colorScheme='primary'
            onClick={handleConfirm}
            dataTest='add-sku-confirm'
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
