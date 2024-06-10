import type {
  Props,
  ControlProps,
  MenuPlacement,
  SelectInstance,
  ClassNamesConfig,
  ClearIndicatorProps,
} from 'react-select';

import ReactSelect from 'react-select';
import { useMemo, forwardRef, useCallback } from 'react';

import merge from 'lodash/merge';
import { match } from 'ts-pattern';
import { twMerge } from 'tailwind-merge';

import { cn } from '@ui/utils/cn';
import { Delete } from '@ui/media/icons/Delete';
import { inputVariants } from '@ui/form/Input/Input';

type Size = 'xs' | 'sm' | 'md' | 'lg';
// Exhaustively typing this Props interface does not offer any benefit at this moment
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface SelectProps extends Props<any, any, any> {
  size?: Size;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
}

export const Select = forwardRef<SelectInstance, SelectProps>(
  (
    {
      isReadOnly,
      leftElement,
      size = 'md',
      components: _components,
      classNames,
      ...rest
    },
    ref,
  ) => {
    const Control = useCallback(
      ({ children, innerRef, innerProps }: ControlProps) => {
        const sizeClass = {
          xs: 'min-h-4',
          sm: 'min-h-7',
          md: 'min-h-8',
          lg: 'min-h-8',
        }[size];

        return (
          <div
            ref={innerRef}
            className={`flex w-full items-center group ${sizeClass}`}
            {...innerProps}
          >
            {leftElement}
            {children}
          </div>
        );
      },
      [leftElement, size],
    );

    const ClearIndicator = useCallback(
      ({ innerProps }: ClearIndicatorProps) => {
        const iconSize = {
          xs: 'size-3',
          sm: 'size-3',
          md: 'size-4',
          lg: 'size-5',
        }[size];

        const wrapperSize = {
          xs: 'size-5',
          sm: 'size-7',
          md: 'size-8',
          lg: 'size-8',
        }[size];

        const { className, ...restInnerProps } = innerProps;

        return (
          <div
            className={cn(
              'flex rounded-md items-center justify-center bg-transparent hover:bg-gray-100',
              wrapperSize,
            )}
            {...restInnerProps}
          >
            <Delete
              className={cn(
                'text-transparent group-hover:text-gray-700',
                iconSize,
              )}
            />
          </div>
        );
      },
      [size],
    );

    const components = useMemo(
      () => ({
        Control,
        ..._components,
        ClearIndicator,
        DropdownIndicator: () => null,
      }),
      [Control, _components],
    );
    const defaultClassNames = useMemo(
      () => merge(getDefaultClassNames({ size, isReadOnly }), classNames),
      [size, isReadOnly, classNames],
    );

    return (
      <ReactSelect
        unstyled
        ref={ref}
        components={components}
        tabSelectsValue={false}
        {...rest}
        classNames={defaultClassNames}
      />
    );
  },
);

const getDefaultClassNames = ({
  size,
  isReadOnly,
}: Pick<SelectProps, 'size' | 'isReadOnly'>): ClassNamesConfig => ({
  container: ({ isFocused }) =>
    getContainerClassNames(undefined, { size, isFocused, isReadOnly }),
  menu: ({ menuPlacement }) =>
    getMenuClassNames(menuPlacement)(
      match(size)
        .with('xs', () => 'text-sm')
        .with('sm', () => 'text-sm')
        .with('md', () => 'text-md')
        .with('lg', () => 'text-lg')
        .otherwise(() => ''),
    ),

  menuList: () =>
    'p-2 max-h-[300px] border border-gray-200 bg-white outline-offset-[2px] outline-[2px] rounded-lg shadow-lg overflow-y-auto overscroll-auto',
  option: ({ isFocused, isSelected }) =>
    cn(
      'my-[2px] px-3 py-1.5 rounded-md text-gray-700 truncate transition ease-in-out delay-50 hover:bg-gray-50',
      isSelected && 'bg-gray-50 font-medium leading-normal',
      isFocused && 'ring-2 ring-gray-100',
    ),
  placeholder: () => 'text-gray-400',
  multiValue: () =>
    'flex items-center h-6 bg-gray-50 rounded-md pl-2 pr-1 ml-0 mr-1 mb-1 border border-gray-200',
  multiValueLabel: () => 'text-gray-500 text-sm mr-1',
  multiValueRemove: () => 'cursor-pointer *:size-5 *:text-gray-500',
  groupHeading: () => 'text-gray-400 text-sm px-3 py-1.5 font-normal uppercase',
});

export const getMultiValueClassNames = (className?: string) => {
  const defaultStyle =
    'flex items-center h-6 bg-gray-50 rounded-md pl-2 pr-1 ml-0 mr-1 mb-1 border border-gray-200';

  return twMerge(defaultStyle, className);
};
export const getMenuClassNames =
  (menuPlacement: MenuPlacement) => (className?: string) => {
    const defaultStyle = cn(
      menuPlacement === 'top'
        ? 'mb-2 animate-slideDownAndFade'
        : 'mt-2 animate-slideUpAndFade',
    );

    return twMerge(defaultStyle, className);
  };
export const getMenuListClassNames = (className?: string) => {
  const defaultStyle =
    'p-2 max-h-[300px] border border-gray-200 bg-white outline-offset-[2px] outline-[2px] rounded-lg shadow-lg overflow-y-auto overscroll-auto';

  return twMerge(defaultStyle, className);
};

export const getMultiValueLabelClassNames = (className?: string) => {
  const defaultStyle = 'text-gray-500 text-sm mr-1';

  return twMerge(defaultStyle, className);
};

export const getContainerClassNames = (
  className?: string,
  props?: {
    size?: Size;
    isFocused?: boolean;
    isReadOnly?: boolean;
  },
) => {
  const defaultStyle = inputVariants({
    variant: 'flushed',
    size: props?.size,
    className: cn(
      'flex items-center cursor-pointer overflow-visible',
      props?.isReadOnly && 'pointer-events-none',
      props?.isFocused && 'border-primary-500',
    ),
  });

  return twMerge(defaultStyle, className);
};
