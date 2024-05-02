import React from 'react';

import { observer } from 'mobx-react-lite';
import { utcToZonedTime } from 'date-fns-tz';

import { cn } from '@ui/utils/cn';
import { BilledType } from '@graphql/types';
import { Delete } from '@ui/media/icons/Delete';
import { DateTimeUtils } from '@spaces/utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
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
  service: ServiceLineItemStore;
}

const billedTypeOptions: SelectOption<BilledType>[] = [
  { label: 'once', value: BilledType.Once },
  { label: 'month', value: BilledType.Monthly },
  { label: 'quarter', value: BilledType.Quarterly },
  { label: 'year', value: BilledType.Annually },
];
const billedTypeLabel: Record<
  Exclude<BilledType, BilledType.None | BilledType.Usage>,
  string
> = {
  [BilledType.Once]: 'once',
  [BilledType.Monthly]: 'month',
  [BilledType.Quarterly]: 'quarter',
  [BilledType.Annually]: 'year',
};
const formSelectClassNames =
  'text-sm inline min-h-1 max-h-3 border-none hover:border-none focus:border-none w-fit ml-1 mt-0 underline  min-w-[max-content]';

const inputClasses =
  'text-sm min-w-2.5 min-h-0 max-h-4 text-inherit underline hover:border-none focus:border-none border-none';
const deleteButtonClasses =
  'border-none bg-transparent shadow-none text-gray-400 px-3 py-2 -mx-3 absolute -right-7 top-0 bottom-0 invisible group-hover:visible hover:bg-transparent';
export const ServiceItem: React.FC<ServiceItemProps> = observer(
  ({ service, isEnded, currency }) => {
    const highlightVersion =
      service.serviceLineItem?.frontendMetadata?.shapeVariant;
    const bgColor = service.serviceLineItem?.frontendMetadata?.color;
    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    return (
      <React.Fragment>
        {service?.serviceLineItem?.isNew &&
        !service.serviceLineItem.isDeleted &&
        !service.serviceLineItem.closedVersion ? (
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
                  onChange={(e) => service.updateQuantity(e.target.value)}
                  size='xs'
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
              <span className='relative z-[2] mx-1'>x</span>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('price') ? bgColor : undefined
                }
              >
                {sliCurrencySymbol}
                <ResizableInput
                  value={service.serviceLineItem?.price}
                  onChange={(e) => service.updatePrice(e.target.value)}
                  size='xs'
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
              <span className='relative z-[2] ml-1'>/</span>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('billingCycle') ? bgColor : undefined
                }
              >
                <Select
                  className={formSelectClassNames}
                  isClearable={false}
                  placeholder='Billed type'
                  value={service.billingValue}
                  onChange={(e) => service.updateBilledType(e.value)}
                  options={billedTypeOptions}
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
              <DatePickerUnderline2
                value={service.serviceLineItem.serviceStarted}
                onChange={(e) => service.updateStartDate(e)}
              />
            </Highlighter>

            <IconButton
              aria-label={'Delete version'}
              icon={<Delete className='text-inherit' />}
              variant='outline'
              size='xs'
              onClick={() => service.setIsDeleted(true)}
              className={deleteButtonClasses}
            />
          </div>
        ) : (
          <div
            className={cn(
              'flex items-baseline justify-between group text-gray-500 relative',
              {
                'text-gray-400': isEnded,
                'line-through text-gray-400 hover:text-gray-400':
                  service.serviceLineItem?.closedVersion ||
                  service.serviceLineItem?.isDeleted,
              },
            )}
          >
            <div className='flex items-baseline'>
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
                      BilledType.None | BilledType.Usage
                    >
                  ]
                }
              </span>
              <span className='ml-1'>
                • {service?.serviceLineItem?.tax?.taxRate}% VAT
              </span>
            </div>
            <div className='ml-1'>
              {service?.serviceLineItem?.serviceStarted &&
                DateTimeUtils.format(
                  utcToZonedTime(
                    service.serviceLineItem.serviceStarted,
                    'UTC',
                  ).toString(),
                  DateTimeUtils.dateWithShortYear,
                )}
            </div>
            {!service.serviceLineItem?.closedVersion &&
              service.serviceLineItem?.isNew && (
                <IconButton
                  aria-label={'Restore version'}
                  icon={<FlipBackward className='text-inherit' />}
                  variant='outline'
                  size='xs'
                  onClick={() => service.setIsDeleted(false)}
                  className={deleteButtonClasses}
                />
              )}
          </div>
        )}
      </React.Fragment>
    );
  },
);

export default ServiceItem;
