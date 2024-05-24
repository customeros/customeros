import React, { useMemo } from 'react';

import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Delete } from '@ui/media/icons/Delete';
import { DateTimeUtils } from '@spaces/utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
import { BilledType, ContractStatus } from '@graphql/types';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencySymbol } from '@shared/util/currencyOptions';
import { ResizableInput } from '@ui/form/Input/ResizableInput';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2';
import {
  Select,
  getMenuClassNames,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select';

import { Highlighter } from './highlighters';
import ServiceLineItemStore from '../../stores/Service.store';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  isModification?: boolean;
  service: ServiceLineItemStore;
<<<<<<< Updated upstream
  type: 'subscription' | 'one-time';
  contractStatus?: ContractStatus | null;
=======
  type?: 'subscription' | 'one-time';
>>>>>>> Stashed changes
}

const billedTypeOptions: SelectOption<BilledType>[] = [
  { label: 'month', value: BilledType.Monthly },
  { label: 'quarter', value: BilledType.Quarterly },
  { label: 'year', value: BilledType.Annually },
];
const billedTypeLabel: Record<
  Exclude<BilledType, BilledType.None | BilledType.Usage | BilledType.Once>,
  string
> = {
  [BilledType.Monthly]: 'month',
  [BilledType.Quarterly]: 'quarter',
  [BilledType.Annually]: 'year',
};
const formSelectClassNames =
  'text-sm inline min-h-1 max-h-3 border-none hover:border-none focus:border-none w-fit ml-1 mt-0 underline  min-w-[max-content]';

const inputClasses =
  'text-sm min-w-2.5 min-h-0 max-h-4 text-inherit underline hover:border-none focus:border-none border-none';

const deleteButtonClasses =
  'border-none bg-transparent shadow-none text-gray-400 pr-3 pl-4 py-2 -mx-4 absolute -right-7 top-0 bottom-0 invisible group-hover:visible hover:bg-transparent';

export const ServiceItem: React.FC<ServiceItemProps> = observer(
<<<<<<< Updated upstream
  ({ service, isEnded, currency, isModification, type, contractStatus }) => {
=======
  ({ service, isEnded, currency, type }) => {
>>>>>>> Stashed changes
    const highlightVersion =
      service.serviceLineItem?.frontendMetadata?.shapeVariant;
    const bgColor = service.serviceLineItem?.frontendMetadata?.color;
    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    const billedTypeOpt = useMemo(() => {
      return type === 'subscription'
        ? billedTypeOptions.filter((opt) => opt.value !== BilledType.Once)
        : billedTypeOptions.filter((opt) => opt.value === BilledType.Once);
    }, [type]);

    const showEditView =
      (contractStatus === ContractStatus.Draft &&
        !service?.serviceLineItem?.closedVersion &&
        !service?.serviceLineItem?.isDeleted) ||
      (service?.serviceLineItem?.isNew &&
        !service.serviceLineItem.isDeleted &&
        !service.serviceLineItem.closedVersion);

    return (
      <React.Fragment>
        {showEditView ? (
          <div className='flex items-baseline justify-between group relative py-3 -my-3 text-gray-500 '>
            <div className='flex items-baseline'>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('quantity') ? bgColor : undefined
                }
              >
                <ResizableInput
                  value={service.serviceLineItem?.quantity}
                  onChange={(e) =>
                    !e.target.value?.length
                      ? service.updateQuantity('0')
                      : service.updateQuantity(e.target.value)
                  }
                  size='xs'
                  type='number'
                  min={0}
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
<<<<<<< Updated upstream
              <span className='relative z-[2] mr-0.5'>x</span>
=======
              <span className='relative z-[2] mx-1'>×</span>
>>>>>>> Stashed changes
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('price') ? bgColor : undefined
                }
              >
                {sliCurrencySymbol}
                <ResizableInput
                  value={service.serviceLineItem?.price}
                  onChange={(e) =>
                    !e.target.value?.length
                      ? service.updatePrice('0')
                      : service.updatePrice(e.target.value)
                  }
                  size='xs'
                  type='number'
                  min={0}
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
<<<<<<< Updated upstream
=======
              {type !== 'one-time' && (
                <span className='relative z-[2] ml-1'>/</span>
              )}
>>>>>>> Stashed changes
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('billingCycle') ? bgColor : undefined
                }
              >
<<<<<<< Updated upstream
                {isModification && contractStatus !== ContractStatus.Draft ? (
                  <span className='text-gray-700'>
                    <span className='mr-0.5'>/</span>
                    {
                      billedTypeLabel[
                        service?.serviceLineItem?.billingCycle as Exclude<
                          BilledType,
                          BilledType.None | BilledType.Usage
                        >
                      ]
                    }
                  </span>
                ) : (
=======
                {type !== 'one-time' && (
>>>>>>> Stashed changes
                  <Select
                    className={formSelectClassNames}
                    isClearable={false}
                    placeholder='Billed type'
                    value={service.billingValue}
                    onChange={(e) => service.updateBilledType(e.value)}
<<<<<<< Updated upstream
                    options={billedTypeOpt}
=======
                    options={billedTypeOptions}
>>>>>>> Stashed changes
                    menuPosition='absolute'
                    classNames={{
                      container: () =>
                        getContainerClassNames(
                          'text-inherit text-base hover:text-gray-500 focus:text-gray-500 min-w-fit w-max-content ml-0',
                          'xs',
                        ),
                      menuList: () => getMenuListClassNames('min-w-[100px]'),
                      menu: ({ menuPlacement }) =>
                        getMenuClassNames(menuPlacement)('!z-[11]'),
                    }}
                    size='xs'
                  />
                )}
              </Highlighter>
              <span className='relative z-[2] mx-1'>•</span>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('taxRate') ? bgColor : undefined
                }
              >
                <ResizableInput
                  value={service.serviceLineItem?.tax?.taxRate}
                  onChange={(e) =>
                    service.updateTaxRate(parseFloat(e.target.value))
                  }
                  size='xs'
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
              <span className='whitespace-nowrap relative z-[2] mx-1'>
                % VAT
              </span>
            </div>

            <Highlighter
              highlightVersion={highlightVersion}
              backgroundColor={
                service.isFieldRevised('serviceStarted') ? bgColor : undefined
              }
            >
              {/* <span className='mr-1'>Current since </span> */}
              <DatePickerUnderline2
                value={service?.serviceLineItem?.serviceStarted}
                minDate={service?.serviceLineItem?.nextBilling ?? undefined}
                onChange={(e) => service.updateStartDate(e)}
              />
            </Highlighter>

            {contractStatus !== ContractStatus.Draft && (
              <IconButton
                aria-label={'Delete version'}
                icon={<Delete className='text-inherit' />}
                variant='outline'
                size='xs'
                onClick={() => {
                  service.setIsDeleted(true);
                }}
                className={deleteButtonClasses}
              />
            )}
          </div>
        ) : (
          <div
            className={cn(
              'flex items-baseline justify-between group text-gray-700 relative',
              {
                'text-gray-400': isEnded,
                'line-through text-gray-400 hover:text-gray-400':
                  service.serviceLineItem?.isDeleted,
              },
            )}
          >
            <div className='flex items-baseline text-gray-700'>
              <span>
                {service?.serviceLineItem?.quantity} x {sliCurrencySymbol}
                {service?.serviceLineItem?.price}
              </span>
              <span>
                /{' '}
                {
                  billedTypeLabel[
                    service?.serviceLineItem?.billingCycle as Exclude<
                      BilledType,
                      BilledType.None | BilledType.Usage | BilledType.Once
                    >
                  ]
                }{' '}
              </span>
              <span className='ml-1 text-gray-700'>
                • {service?.serviceLineItem?.tax?.taxRate}% VAT
              </span>
            </div>
            <div className='ml-1 text-gray-700'>
              {service?.serviceLineItem?.serviceStarted &&
                DateTimeUtils.format(
                  toZonedTime(
                    service.serviceLineItem.serviceStarted,
                    'UTC',
                  ).toString(),
                  DateTimeUtils.dateWithShortYear,
                )}
            </div>
            {service.serviceLineItem?.isNew ||
              (contractStatus === ContractStatus.Draft &&
                service.serviceLineItem?.isDeleted && (
                  <IconButton
                    aria-label={'Restore version'}
                    icon={<FlipBackward className='text-inherit' />}
                    variant='outline'
                    size='xs'
                    onClick={() => service.setIsDeleted(false)}
                    className={deleteButtonClasses}
                  />
                ))}
          </div>
        )}
      </React.Fragment>
    );
  },
);

export default ServiceItem;
