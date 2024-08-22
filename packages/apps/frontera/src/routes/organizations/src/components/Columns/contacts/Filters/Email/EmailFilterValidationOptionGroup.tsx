import React, { MouseEvent } from 'react';

import { CheckedState } from '@radix-ui/react-checkbox';

import { IconButton } from '@ui/form/IconButton';
import { InfoCircle } from '@ui/media/icons/InfoCircle.tsx';
import { Checkbox, CheckMinus } from '@ui/form/Checkbox/Checkbox';
import {
  CollapsibleRoot,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@ui/transitions/Collapse/Collapse.tsx';

import {
  CategoryHeaderLabel,
  DeliverabilityStatus,
  EmailVerificationStatus,
} from './utils.ts';

interface CheckboxOption {
  label: string;
  disabled?: boolean;
  value: EmailVerificationStatus;
}

interface CheckboxGroupProps {
  options: CheckboxOption[];
  isCategoryChecked: boolean;
  category: DeliverabilityStatus;
  onToggleCategory: (category: DeliverabilityStatus) => void;
  isOptionChecked: (value: EmailVerificationStatus) => boolean;
  onToggleOption: (
    value: EmailVerificationStatus,
    checked?: CheckedState,
  ) => void;
  onOpenInfoModal: (
    e: MouseEvent<HTMLButtonElement>,
    value: EmailVerificationStatus,
  ) => void;
}

export const EmailFilterValidationOptionGroup: React.FC<CheckboxGroupProps> = ({
  category,
  options,
  onToggleCategory,
  onToggleOption,
  isOptionChecked,
  isCategoryChecked,
  onOpenInfoModal,
}) => (
  <CollapsibleRoot open={isCategoryChecked} className='flex flex-col w-full'>
    <div className='flex justify-between w-full items-center'>
      <CollapsibleTrigger asChild={false}>
        <Checkbox
          icon={<CheckMinus />}
          isChecked={isCategoryChecked}
          onChange={() => onToggleCategory(category)}
        >
          <p className='text-sm'>{CategoryHeaderLabel[category]}</p>
        </Checkbox>
      </CollapsibleTrigger>
    </div>
    <CollapsibleContent>
      <div className='flex flex-col w-full gap-2 ml-6 mt-2'>
        {options.map((option) => (
          <div className='group' key={option.value}>
            <Checkbox
              disabled={option?.disabled}
              isChecked={isOptionChecked(option.value)}
              onChange={(checked) => onToggleOption(option.value, checked)}
            >
              <div className='flex '>
                <p className='text-sm'>{option.label}</p>

                <IconButton
                  size='xxs'
                  variant='ghost'
                  icon={<InfoCircle />}
                  aria-label='More info'
                  onClick={(e) => onOpenInfoModal(e, option.value)}
                  className='opacity-0 group-hover:opacity-100 transition-opacity bg-transparent'
                />
              </div>
            </Checkbox>
          </div>
        ))}
      </div>
    </CollapsibleContent>
  </CollapsibleRoot>
);
