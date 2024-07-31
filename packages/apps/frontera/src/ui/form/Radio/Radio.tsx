import React, { ElementRef, forwardRef, ComponentPropsWithoutRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixRadioGroup from '@radix-ui/react-radio-group';

interface RadioProps {
  children?: React.ReactNode;
}

export const RadioGroup = forwardRef<
  ElementRef<typeof RadixRadioGroup.Root> & RadioProps,
  ComponentPropsWithoutRef<typeof RadixRadioGroup.Root>
>(({ className, children, ...props }, ref) => {
  return (
    <RadixRadioGroup.Root
      ref={ref}
      aria-label='View density'
      className={twMerge('flex flex-col gap-2.5', className)}
      {...props}
    >
      {children}
    </RadixRadioGroup.Root>
  );
});

interface RadioItemProps {
  children?: React.ReactNode;
}

export const Radio = forwardRef<
  ElementRef<typeof RadixRadioGroup.Item> & RadioItemProps,
  ComponentPropsWithoutRef<typeof RadixRadioGroup.Item>
>(({ className, children, ...props }, ref) => {
  return (
    <div className='flex space-x-2 items-center'>
      <RadixRadioGroup.Item
        ref={ref}
        className={twMerge(
          'bg-white size-4 rounded-full border border-solid border-gray-300 hover:border-primary-500 hover:bg-primary-50 focus:ring-4 focus:ring-primary-50 data-[state=checked]:bg-primary-50 data-[state=checked]:border-primary-500 data-[disabled]:border-gray-300 data-[disabled]:bg-gray-100 data-[disabled]:cursor-not-allowed  outline-none cursor-pointer',
          className,
        )}
        {...props}
      >
        <RadixRadioGroup.Indicator className='flex items-center justify-center w-full h-full relative data-[disabled]:after:bg-gray-300 data-[state=checked]:after:bg-primary-600 data-[state=checked]:after:rounded-full data-[state=checked]:after:size-2 data-[state=checked]:after:block data-[state=checked]:after:content-[""]' />
      </RadixRadioGroup.Item>
      {children}
    </div>
  );
});
