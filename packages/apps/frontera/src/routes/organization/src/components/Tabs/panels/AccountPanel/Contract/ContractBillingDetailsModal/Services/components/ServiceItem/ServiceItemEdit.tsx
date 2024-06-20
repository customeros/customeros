import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date.ts';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { toastError } from '@ui/presentation/Toast';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { currencySymbol } from '@shared/util/currencyOptions.ts';
import { ResizableInput } from '@ui/form/Input/ResizableInput.tsx';
import { BilledType, ContractStatus, ServiceLineItem } from '@graphql/types';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2.tsx';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

import { Highlighter } from '../highlighters';

interface ServiceItemProps {
  currency?: string;
  isModification?: boolean;
  service: Store<ServiceLineItem>;

  type: 'subscription' | 'one-time';
  allServices?: Store<ServiceLineItem>[];

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

export const ServiceItemEdit: React.FC<ServiceItemProps> = observer(
  ({
    service,
    allServices,
    currency,
    isModification,
    type,
    contractStatus,
  }) => {
    const highlightVersion = '';
    // service?.value?.frontendMetadata?.shapeVariant;

    // const bgColor = billingEnabled ? 'transparent' : 'transparent';
    // ? service?.value?.frontendMetadata?.color
    // : 'transparent';

    const sliCurrencySymbol = currency ? currencySymbol?.[currency] : '$';

    const isDraft =
      contractStatus &&
      [ContractStatus.Draft, ContractStatus.Scheduled].includes(contractStatus);

    const onChangeServiceStarted = (e: Date | null) => {
      if (!e) return;

      const checkExistingServiceStarted = (date: Date) => {
        return allServices?.some((service) =>
          DateTimeUtils.isSameDay(service?.value?.serviceStarted, `${date}`),
        );
      };

      const findCurrentService = () => {
        if (isDraft) return null;

        return allServices?.find((serviceData) => {
          const serviceStarted = serviceData?.value?.serviceStarted;
          const serviceEnded = serviceData?.value?.serviceEnded;

          return (
            (serviceEnded &&
              DateTimeUtils.isFuture(serviceEnded) &&
              DateTimeUtils.isPast(serviceStarted)) ||
            (!serviceEnded && DateTimeUtils.isPast(serviceStarted))
          );
        })?.value?.serviceStarted;
      };

      const checkIfBeforeCurrentService = (
        date: Date,
        currentService: Date | null,
      ) => {
        return (
          currentService &&
          DateTimeUtils.isBefore(date.toString(), currentService.toString())
        );
      };

      const existingServiceStarted = checkExistingServiceStarted(e);
      const currentService = findCurrentService();
      const isBeforeCurrentService = checkIfBeforeCurrentService(
        e,
        currentService,
      );

      if (isBeforeCurrentService) {
        toastError(
          `Modifications must be effective after the current service`,
          `${service?.value?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      if (existingServiceStarted) {
        toastError(
          `A version with this date already exists`,
          `${service?.value?.metadata?.id}-service-started-date-update-error`,
        );

        return;
      }

      service.update(
        (prev) => ({
          ...prev,
          serviceStarted: e,
        }),
        { mutate: false },
      );
    };

    const updateQuantity = (quantity: string) => {
      service.update((prev) => ({ ...prev, quantity }), { mutate: false });
    };
    const updatePrice = (price: string) => {
      service.update(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-expect-error  we allow undefined during edition but on blur we still enforce value therefore this is false positive
        (prev) => ({ ...prev, price: price ? parseFloat(price) : undefined }),
        { mutate: false },
      );
    };
    const updateTaxRate = (taxRate: string) => {
      service.update(
        (prev) => ({
          ...prev,
          tax: {
            ...prev.tax,
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error we allow undefined during edition but on blur we still enforce value therefore this is false positive
            taxRate: taxRate ? parseFloat(taxRate) : undefined,
          },
        }),
        {
          mutate: false,
        },
      );
    };

    return (
      <div className='flex items-baseline justify-between group relative text-gray-500 '>
        <div className='flex items-baseline'>
          <Highlighter
            highlightVersion={highlightVersion}
            backgroundColor={
              undefined
              // service.isFieldRevised('quantity') ? bgColor : undefined
            }
          >
            <ResizableInput
              value={service?.value?.quantity ?? ''}
              onChange={(e) => updateQuantity(e.target.value ?? '')}
              onBlur={(e) =>
                !e.target.value?.length
                  ? updateQuantity('0')
                  : updateQuantity(e.target.value)
              }
              placeholder='0'
              size='xs'
              type='number'
              min={0}
              className={inputClasses}
              onFocus={(e) => e.target.select()}
            />
          </Highlighter>
          <span className='relative z-[2] mx-1 text-gray-700'>×</span>
          <Highlighter
            highlightVersion={highlightVersion}
            backgroundColor={
              undefined
              // service.isFieldRevised('price') ? bgColor : undefined
            }
          >
            {sliCurrencySymbol}
            <ResizableInput
              value={service?.value?.price}
              onChange={(e) => updatePrice(e.target.value ?? '')}
              onBlur={(e) =>
                !e.target.value?.length
                  ? updatePrice('0')
                  : updatePrice(e.target.value)
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
              undefined
              // service.isFieldRevised('price') ? bgColor : undefined
            }
          >
            {type === 'one-time' ? (
              <span className='text-gray-700'></span>
            ) : !service.value?.metadata?.id ? (
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
                            service?.value?.billingCycle as Exclude<
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
                        service.update((prev) => ({
                          ...prev,
                          billingCycle: option.value,
                        }));
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
                  service.value &&
                  billedTypesLabel(
                    service?.value?.billingCycle.toLocaleLowerCase(),
                  )}
              </p>
            )}
          </Highlighter>
          <span className='relative z-[2] mx-1 text-gray-700'>•</span>
          <Highlighter
            highlightVersion={highlightVersion}
            backgroundColor={
              undefined
              // service.isFieldRevised('taxRate') ? bgColor : undefined
            }
          >
            <ResizableInput
              value={
                !isNaN(service?.value?.tax?.taxRate as number)
                  ? service?.value?.tax.taxRate
                  : ''
              }
              onChange={(e) => updateTaxRate(e.target.value)}
              onBlur={(e) =>
                !e.target.value?.trim()?.length
                  ? updateTaxRate('0')
                  : updateTaxRate(e.target.value)
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
            undefined
            // service.isFieldRevised('serviceStarted') ? bgColor : undefined
          }
        >
          <DatePickerUnderline2
            value={service?.value?.serviceStarted}
            onChange={onChangeServiceStarted}
          />
        </Highlighter>

        <IconButton
          aria-label={'Delete version'}
          icon={<Delete className='text-inherit' />}
          variant='outline'
          size='xs'
          onClick={() => {
            service.update((prev) => ({ ...prev, closed: true }), {
              mutate: false,
            });
          }}
          className={deleteButtonClasses}
        />
      </div>
    );
  },
);
