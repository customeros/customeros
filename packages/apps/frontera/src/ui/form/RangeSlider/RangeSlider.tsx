import React, { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as RadixSlider from '@radix-ui/react-slider';

interface RangeSliderProps extends RadixSlider.SliderProps {
  className?: string;
  children: React.ReactNode;
}

export const RangeSlider = forwardRef<HTMLSpanElement, RangeSliderProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <RadixSlider.Root
        ref={ref}
        className={twMerge(
          'relative flex items-center select-none touch-none w-full h-5',
          className,
        )}
        {...props}
      >
        {children}
      </RadixSlider.Root>
    );
  },
);
interface RangeSliderGenericProps {
  dataTest?: string;
  className?: string;
  children?: React.ReactNode;
}

export const RangeSliderTrack = ({
  children,
  className,
  dataTest,
}: RangeSliderGenericProps) => {
  return (
    <RadixSlider.Track
      data-test={dataTest}
      className={twMerge('relative grow rounded-full', className)}
    >
      {children}
    </RadixSlider.Track>
  );
};

export const RangeSliderFilledTrack = ({
  className,
}: RangeSliderGenericProps) => {
  return (
    <RadixSlider.Range
      className={twMerge('absolute rounded-full', className)}
    />
  );
};

export const RangeSliderThumb = ({ className }: RangeSliderGenericProps) => {
  return (
    <RadixSlider.Thumb
      className={twMerge('block w-5 h-5 bg-white rounded-[10px]', className)}
    />
  );
};
