import ReactSelect from 'react-select';
import { useMemo, forwardRef, useCallback } from 'react';
import type {
  Props,
  ControlProps,
  MenuPlacement,
  SelectInstance,
  ClassNamesConfig,
  ClearIndicatorProps,
} from 'react-select';

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
  dataTest?: string;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
  onKeyDown?: (e: React.KeyboardEvent<HTMLDivElement>) => void;
}

export const Select = forwardRef<SelectInstance, SelectProps>(
  (
    {
      isReadOnly,
      leftElement,
      size = 'md',
      dataTest,
      components: _components,
      classNames,
      onKeyDown,
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
            data-test={dataTest}
          >
            {leftElement}
            {children}
          </div>
        );
      },
      [leftElement, size, dataTest],
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
                'text-transparent group-hover:text-gray-700 ',
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
        onKeyDown={(e) => {
          if (onKeyDown) {
            onKeyDown(e);
            e.stopPropagation();
          }
        }}
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
    getContainerClassNames(undefined, 'flushed', {
      size,
      isFocused,
      isReadOnly,
    }),
  menu: ({ menuPlacement }) => getMenuClassNames(menuPlacement)('', size),
  menuList: () => getMenuListClassNames(),
  option: ({ isFocused, isSelected }) =>
    getOptionClassNames('', { isFocused, isSelected }),
  placeholder: () => 'text-gray-400',
  multiValue: () => getMultiValueClassNames(''),
  multiValueLabel: () => getMultiValueLabelClassNames('', size),
  multiValueRemove: () => getMultiValueRemoveClassNames('', size),
  groupHeading: () => 'text-gray-400 text-sm px-3 py-1.5 font-normal',
  valueContainer: () => 'gap-1 py-0.5 mr-0.5 inline-grid',
});

export const getMultiValueRemoveClassNames = (
  className?: string,
  size?: string,
) => {
  const sizeClass = match(size)
    .with('xs', () => 'size-5 *:size-5')
    .with('sm', () => 'size-5 *:size-5')
    .with('md', () => 'size-6 *:size-6')
    .with('lg', () => 'size-7 *:size-7')
    .otherwise(() => '');

  return twMerge(
    'cursor-pointer text-grayModern-400 mr-0 bg-grayModern-100 rounded-e-md px-0.5 hover:bg-grayModern-200 hover:text-warning-700 transition ease-in-out',
    sizeClass,
    className,
  );
};

export const getMultiValueClassNames = (className?: string) => {
  const defaultStyle = 'border-none mb-0 bg-transparent mr-0 pl-0';

  return twMerge(defaultStyle, className);
};

export const getMenuClassNames =
  (menuPlacement: MenuPlacement) => (className?: string, size?: Size) => {
    const sizes = match(size)
      .with('xs', () => 'text-sm')
      .with('sm', () => 'text-sm')
      .with('md', () => 'text-md')
      .with('lg', () => 'text-lg')
      .otherwise(() => '');

    const defaultStyle = cn(
      menuPlacement === 'top'
        ? 'mb-2 animate-slideDownAndFade'
        : 'mt-2 animate-slideUpAndFade',
    );

    return twMerge(defaultStyle, sizes, className);
  };

export const getMenuListClassNames = (className?: string) => {
  const defaultStyle =
    'p-2 max-h-[300px] border border-gray-200 bg-white outline-offset-[2px] outline-[2px] rounded-lg shadow-lg overflow-y-auto overscroll-auto';

  return twMerge(defaultStyle, className);
};

export const getMultiValueLabelClassNames = (
  className?: string,
  size?: string,
) => {
  const sizeClass = match(size)
    .with('xs', () => 'text-sm')
    .with('sm', () => 'text-sm')
    .with('md', () => 'text-md')
    .with('lg', () => 'text-lg')
    .otherwise(() => '');

  const defaultStyle = cn(
    'bg-grayModern-100 text-gray-700 px-1 mr-0 rounded-s-md hover:bg-grayModern-200 transition ease-in-out',
    sizeClass,
  );

  return twMerge(defaultStyle, className);
};

export const getContainerClassNames = (
  className?: string,
  variant?: 'flushed' | 'unstyled' | 'group' | 'outline',
  props?: {
    size?: Size;
    isFocused?: boolean;
    isReadOnly?: boolean;
  },
) => {
  const defaultStyle = inputVariants({
    variant: variant || 'flushed',
    size: props?.size,
    className: cn(
      'flex items-center cursor-pointer overflow-visible outline-0',
      props?.isReadOnly && 'pointer-events-none',
      props?.isFocused && 'border-primary-500',
    ),
  });

  return twMerge(defaultStyle, className, variant);
};

export const getOptionClassNames = (
  className: string = '',
  props: { isFocused?: boolean; isSelected?: boolean },
) => {
  const { isFocused, isSelected } = props;

  return cn(
    'my-[2px] px-3 py-1 rounded-md text-gray-700 truncate transition ease-in-out delay-50 hover:bg-grayModern-100',
    isSelected && 'bg-gray-50 font-medium leading-normal',
    isFocused && 'bg-grayModern-100',
    className,
  );
};

export { components } from 'react-select';
export type { OptionProps } from 'react-select';
