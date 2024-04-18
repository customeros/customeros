import React, {
  Ref,
  useRef,
  useEffect,
  forwardRef,
  HTMLAttributes,
} from 'react';

import { twMerge } from 'tailwind-merge';

import { cn } from '@ui/utils/cn';

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
    const isFirstRender = useIsFirstRender();
    const wasPreviouslyLoaded = usePrevious(isLoaded);

    if (isLoaded) {
      return (
        <div
          ref={ref}
          {...props}
          className={cn(
            isFirstRender || wasPreviouslyLoaded ? 'none' : defaultClasses,
            className,
          )}
        />
      );
    }

    return (
      <div
        ref={ref}
        {...props}
        className={twMerge(defaultClasses, className)}
      />
    );
  },
);

function useIsFirstRender() {
  const isFirstRender = useRef(true);

  useEffect(() => {
    isFirstRender.current = false;
  }, []);

  return isFirstRender.current;
}

function usePrevious<T>(value: T) {
  const ref = useRef<T | undefined>();

  useEffect(() => {
    ref.current = value;
  }, [value]);

  return ref.current as T;
}
