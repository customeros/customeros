import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { BilledType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface BilledTypeEditFieldProps {
  id: string;
  isModification?: boolean;
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

export const BilledTypeEditField: FC<BilledTypeEditFieldProps> = observer(
  ({ id, isModification }) => {
    const store = useStore();
    const service = store.contractLineItems.value.get(id);
    if (
      !service?.value?.metadata?.id.includes('new') ||
      service?.value?.parentId
    ) {
      return (
        <p className='text-gray-700'>
          /
          {service &&
            service.value &&
            billedTypesLabel(service?.value?.billingCycle.toLocaleLowerCase())}
        </p>
      );
    }

    return (
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
                      BilledType.None | BilledType.Usage | BilledType.Once
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
    );
  },
);
