import { useRef, useState, useEffect, KeyboardEvent } from 'react';

import { P, match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { Currency } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
import { MaskedInput } from '@ui/form/Input/MaskedInput';
import { currencySymbol } from '@shared/util/currencyOptions';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber.ts';

interface ArrEstimateProps {
  opportunityId: string;
}

export const ArrEstimateCell = observer(
  ({ opportunityId }: ArrEstimateProps) => {
    const store = useStore();
    const valueInputRef = useRef<HTMLInputElement | null>(null);

    const opportunity = store.opportunities.value.get(opportunityId);
    const defaultValue = opportunity?.value?.maxAmount ?? 0;
    const [value, setValue] = useState(defaultValue.toString());
    const [isHovered, setIsHovered] = useState(false);

    const [isEdit, setIsEdit] = useState(false);
    const ref = useRef(null);

    useOutsideClick({
      ref: ref,
      handler: () => {
        setIsEdit(false);
      },
    });

    useEffect(() => {
      if (isHovered && isEdit) {
        valueInputRef.current?.focus();
      }
    }, [isHovered, isEdit]);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit]);

    const handleEscape = (e: KeyboardEvent<HTMLDivElement>) => {
      if (e.key === 'Escape' || e.key === 'Enter') {
        valueInputRef?.current?.blur();
        setIsEdit(false);
      }
    };

    const handleAccept = (unmaskedValue: string) => {
      opportunity?.update(
        (value) => {
          value.maxAmount = unmaskedValue ? parseFloat(unmaskedValue) : 0;

          return value;
        },
        { mutate: false },
      );
    };

    const handleBlur = () => {
      opportunity?.saveProperty('maxAmount');
    };

    const defaultCurrency = match(store.settings.tenant.value?.baseCurrency)
      .with(P.nullish, () => Currency.Usd)
      .with(P.string, (str) => (str.length === 3 ? str : Currency.Usd))
      .otherwise((tenantCurrency) => tenantCurrency);

    const symbol = match(opportunity?.value?.currency)
      .with(P.nullish, () => currencySymbol[defaultCurrency])
      .otherwise(
        (currency) =>
          currencySymbol[currency] ?? currencySymbol[defaultCurrency],
      );

    useEffect(() => {
      setValue(defaultValue.toString());
    }, [defaultValue]);

    return (
      <div
        ref={ref}
        onKeyDown={handleEscape}
        className='flex justify-between'
        onDoubleClick={() => setIsEdit(true)}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <div className='flex ' style={{ width: `calc(100% - 1rem)` }}>
          {!isEdit && !value && <p className='text-gray-400'>No estimate</p>}
          {!isEdit && value && (
            <p className='overflow-ellipsis overflow-hidden'>
              {formatCurrency(
                opportunity?.value?.maxAmount || 0,
                2,
                opportunity?.value?.currency || defaultCurrency,
              )}
            </p>
          )}
          {isEdit && (
            <MaskedInput
              size='xs'
              value={value}
              variant='unstyled'
              onBlur={handleBlur}
              mask={`${symbol}num`}
              placeholder='ARR estimate'
              defaultValue={defaultValue.toString()}
              onClick={(e) => (e.target as HTMLInputElement).select()}
              onAccept={(v, instance) => {
                setValue(v);
                handleAccept(instance._unmaskedValue);
              }}
              blocks={{
                num: {
                  mask: Number,
                  scale: 0,
                  lazy: false,
                  min: 0,
                  placeholderChar: '#',
                  thousandsSeparator: ',',
                  normalizeZeros: true,
                  padFractionalZeros: true,
                  autofix: true,
                },
              }}
            />
          )}
          {isHovered && !isEdit && (
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='edit'
              className='ml-3 rounded-[5px]'
              onClick={() => setIsEdit(!isEdit)}
              icon={<Edit03 className='text-gray-500' />}
            />
          )}
        </div>
      </div>
    );
  },
);
