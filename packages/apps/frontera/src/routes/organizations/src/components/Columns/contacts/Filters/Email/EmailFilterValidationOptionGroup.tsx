import React, { MouseEvent } from 'react';

import { CheckedState } from '@radix-ui/react-checkbox';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { EmailDeliverable } from '@graphql/types';
import { Checkbox } from '@ui/form/Checkbox/Checkbox';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import {
  CollapsibleRoot,
  CollapsibleContent,
} from '@ui/transitions/Collapse/Collapse';

import { CategoryHeaderLabel, EmailVerificationStatus } from './utils';

interface CheckboxOption {
  label: string;
  disabled?: boolean;
  value: EmailVerificationStatus;
}

interface CheckboxGroupProps {
  options: CheckboxOption[];
  category: EmailDeliverable;
  isCategoryChecked: boolean | CheckedState;
  onToggleCategory: (category: EmailDeliverable) => void;
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
  <CollapsibleRoot open={!!isCategoryChecked} className='flex flex-col w-full'>
    <div className='flex justify-between w-full items-center'>
      <Checkbox
        isChecked={isCategoryChecked}
        onChange={() => onToggleCategory(category)}
      >
        <p className='text-sm'>{CategoryHeaderLabel[category]}</p>
      </Checkbox>
    </div>
    <CollapsibleContent>
      <div className='flex flex-col w-full gap-2 ml-6 mt-2'>
        {options.map((option) => (
          <div className='group' key={option.value}>
            <div className='flex'>
              <Checkbox
                disabled={option?.disabled}
                isChecked={isOptionChecked(option.value)}
                onChange={(checked) => onToggleOption(option.value, checked)}
              >
                <p
                  className={cn('text-sm', {
                    'text-gray-400': option?.disabled,
                  })}
                >
                  {option.label}
                </p>
              </Checkbox>
              <IconButton
                size='xxs'
                variant='ghost'
                icon={<InfoCircle />}
                aria-label='More info'
                onClick={(e) => onOpenInfoModal(e, option.value)}
                className='opacity-0 group-hover:opacity-100 transition-opacity bg-transparent hover:bg-transparent'
              />
            </div>
          </div>
        ))}
      </div>
    </CollapsibleContent>
  </CollapsibleRoot>
);
