import React, { Ref, forwardRef, HTMLAttributes } from 'react';

import { twMerge } from 'tailwind-merge';

interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  isLoaded?: boolean;
  className?: string;
}

const defaultClasses = 'animate-pulse bg-gray-300 rounded-sm';

export const Skeleton = forwardRef(
  (
    { className, isLoaded, ...props }: SkeletonProps,
    ref: Ref<HTMLDivElement>,
  ) => {
    return (
      <div
        ref={ref}
        {...props}
        className={twMerge(defaultClasses, className)}
      />
    );
  },
);
