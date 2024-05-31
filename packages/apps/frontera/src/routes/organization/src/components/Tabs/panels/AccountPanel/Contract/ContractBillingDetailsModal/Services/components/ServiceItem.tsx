import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Delete } from '@ui/media/icons/Delete';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { SelectOption } from '@shared/types/SelectOptions';
import { BilledType, ContractStatus } from '@graphql/types';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { currencySymbol } from '@shared/util/currencyOptions';
import { ResizableInput } from '@ui/form/Input/ResizableInput';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2';

import { Highlighter } from './highlighters';
import ServiceLineItemStore from '../../stores/Service.store';

interface ServiceItemProps {
  isEnded?: boolean;
  currency?: string;
  isModification?: boolean;
  service: ServiceLineItemStore;
  type: 'subscription' | 'one-time';
  contractStatus?: ContractStatus | null;
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

const billedTypesLabel = (label: string) => {
  switch (label) {
    case 'monthly':
      return 'month';
    case 'quarterly':
      return 'quarter';
    case 'annually':
      return 'year';
    default:
      return '';
  }
};

const inputClasses =
  'text-sm min-w-2.5 min-h-0 max-h-4 text-inherit underline hover:border-none focus:border-none border-none';

const deleteButtonClasses =
  'border-none bg-transparent shadow-none text-gray-400 pr-3 pl-4 py-2 -mx-4 absolute -right-7 top-0 bottom-0 invisible group-hover:visible hover:bg-transparent';

export const ServiceItem: React.FC<ServiceItemProps> = observer(
  ({ service, isEnded, currency, isModification, type, contractStatus }) => {
    const highlightVersion =
      service.serviceLineItem?.frontendMetadata?.shapeVariant;

    const bgColor = service.serviceLineItem?.frontendMetadata?.color;

    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    const isFutureVersion =
      service?.serviceLineItem?.serviceStarted &&
      DateTimeUtils.isFuture(service?.serviceLineItem?.serviceStarted);

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const showEditView =
      (isDraft && !service.serviceLineItem?.isDeleted) ||
      (isFutureVersion && !service.serviceLineItem?.isDeleted) ||
      (service?.serviceLineItem?.isNew &&
        !service.serviceLineItem.isDeleted &&
        !service.serviceLineItem.closedVersion);

    const isCurrentVersion =
      (service?.serviceLineItem?.serviceEnded &&
        DateTimeUtils.isFuture(service?.serviceLineItem?.serviceEnded) &&
        DateTimeUtils.isPast(service?.serviceLineItem?.serviceStarted)) ||
      (!service?.serviceLineItem?.serviceEnded &&
        DateTimeUtils.isPast(service?.serviceLineItem?.serviceStarted));

    return (
      <>
        {showEditView ? (
          <div className='flex items-baseline justify-between group relative text-gray-500 '>
            <div className='flex items-baseline'>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('quantity') ? bgColor : undefined
                }
              >
                <ResizableInput
                  value={service.serviceLineItem?.quantity ?? ''}
                  onChange={(e) => service.updateQuantity(e.target.value ?? '')}
                  onBlur={(e) =>
                    !e.target.value?.length
                      ? service.updateQuantity('0')
                      : service.updateQuantity(e.target.value)
                  }
                  placeholder='0'
                  size='xs'
                  type='number'
                  min={0}
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
              <span className='relative z-[2] mx-0.5 text-gray-700'>×</span>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('price') ? bgColor : undefined
                }
              >
                {sliCurrencySymbol}
                <ResizableInput
                  value={service.serviceLineItem?.price}
                  onChange={(e) => service.updatePrice(e.target.value ?? '')}
                  onBlur={(e) =>
                    !e.target.value?.length
                      ? service.updatePrice('0')
                      : service.updatePrice(e.target.value)
                  }
                  size='xs'
                  placeholder='0'
                  type='number'
                  min={0}
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                />
              </Highlighter>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('billingCycle') ? bgColor : undefined
                }
              >
                {isModification || type === 'one-time' ? (
                  <span className='text-gray-700'></span>
                ) : service.isNewlyAdded ? (
                  <Menu>
                    <MenuButton>
                      {isModification ? (
                        <span className='text-gray-700'>
                          <span className='mr-0.5 underline'>/</span>
                        </span>
                      ) : (
                        <span className='text-gray-700'>
                          <span className='mr-0.5'>/</span>
                          <span className='underline text-gray-500'>
                            {
                              billedTypeLabel[
                                service?.serviceLineItem
                                  ?.billingCycle as Exclude<
                                  BilledType,
                                  | BilledType.None
                                  | BilledType.Usage
                                  | BilledType.Once
                                >
                              ]
                            }
                          </span>
                        </span>
                      )}
                    </MenuButton>

                    <MenuList className='min-w-[100px]'>
                      {billedTypeOptions.map((option) => (
                        <MenuItem
                          key={option.value}
                          onClick={() => {
                            service.updateBilledType(option.value);
                          }}
                        >
                          {option.label}
                        </MenuItem>
                      ))}
                    </MenuList>
                  </Menu>
                ) : (
                  <p className='text-gray-700'>
                    /
                    {service &&
                      service.serviceLineItem &&
                      billedTypesLabel(
                        service.serviceLineItem?.billingCycle.toLocaleLowerCase(),
                      )}
                  </p>
                )}
              </Highlighter>
              <span className='relative z-[2] mx-1 text-gray-700'>•</span>
              <Highlighter
                highlightVersion={highlightVersion}
                backgroundColor={
                  service.isFieldRevised('taxRate') ? bgColor : undefined
                }
              >
                <ResizableInput
                  value={service.serviceLineItem?.tax?.taxRate ?? ''}
                  onChange={(e) => service.updateTaxRate(e.target.value ?? '')}
                  onBlur={(e) =>
                    !e.target.value?.length
                      ? service.updateTaxRate('0')
                      : service.updateTaxRate(e.target.value)
                  }
                  placeholder='0'
                  size='xs'
                  className={inputClasses}
                  onFocus={(e) => e.target.select()}
                  min={0}
                />
              </Highlighter>
              <span className='whitespace-nowrap relative z-[2] mx-1 text-gray-700'>
                % VAT
              </span>
            </div>

            <Highlighter
              highlightVersion={highlightVersion}
              backgroundColor={
                service.isFieldRevised('serviceStarted') ? bgColor : undefined
              }
            >
              <Tooltip
                label={
                  service?.serviceLineItem?.nextBilling && !isDraft
                    ? `Please ensure that the service start date is after the start dates of all prior versions of this service`
                    : ''
                }
              >
                <span>
                  <DatePickerUnderline2
                    value={service?.serviceLineItem?.serviceStarted}
                    onChange={(e) => service.updateStartDate(e)}
                  />
                </span>
              </Tooltip>
            </Highlighter>

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
            <div className='flex items-baseline text-inherit'>
              <span>
                {service?.serviceLineItem?.quantity}
                <span className='relative z-[2] mx-1'>×</span>

                {sliCurrencySymbol}
                {service?.serviceLineItem?.price}
              </span>
              {type !== 'one-time' && <span>/ </span>}
              <span>
                {' '}
                {
                  billedTypeLabel[
                    service?.serviceLineItem?.billingCycle as Exclude<
                      BilledType,
                      BilledType.None | BilledType.Usage | BilledType.Once
                    >
                  ]
                }{' '}
              </span>

              <span className='ml-1 text-inherit'>
                • {service?.serviceLineItem?.tax?.taxRate}% VAT
              </span>
            </div>

            <div className='ml-1 text-inherit'>
              {isCurrentVersion && 'Current since '}

              {service?.serviceLineItem?.serviceStarted &&
                DateTimeUtils.format(
                  toZonedTime(
                    service.serviceLineItem.serviceStarted,
                    'UTC',
                  ).toString(),
                  DateTimeUtils.dateWithShortYear,
                )}
            </div>
            {(service.serviceLineItem?.isNew || isDraft || isFutureVersion) &&
              service.serviceLineItem?.isDeleted && (
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
      </>
    );
  },
);

export default ServiceItem;
