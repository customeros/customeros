import { useMemo, forwardRef, useCallback, ComponentType } from 'react';
import {
  OptionProps,
  ControlProps,
  MenuListProps,
  MenuPlacement,
  SelectInstance,
  MultiValueProps,
  ClassNamesConfig,
  MultiValueGenericProps,
  components as createComponent,
} from 'react-select';

import merge from 'lodash/merge';
import { match } from 'ts-pattern';
import { twMerge } from 'tailwind-merge';
import AsyncCreatableSelect, {
  AsyncCreatableProps,
} from 'react-select/async-creatable';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

import { SelectProps } from '../Select';
import { inputVariants } from '../Input';

type Size = 'xs' | 'sm' | 'md' | 'lg';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name?: string;
  formId?: string;
  withTooltip?: boolean;
  leftElement?: React.ReactNode;
  size?: 'xs' | 'sm' | 'md' | 'lg';
  navigateAfterAddingToContract?: boolean;
  removeValue?: (value: SelectOption) => void;
  optionAction?: (data: string) => JSX.Element;
  Option?: ComponentType<OptionProps<SelectOption>>;
  MultiValue?: ComponentType<MultiValueProps<SelectOption>>;
}

export const CreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  (
    {
      size = 'md',
      name,
      formId,
      leftElement,
      components: _components,
      classNames,
      ...props
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
            className={`flex w-full items-start group ${sizeClass}`}
            {...innerProps}
          >
            {leftElement}
            {children}
          </div>
        );
      },
      [size],
    );

    const Option = useCallback(
      ({ data, isFocused, innerRef, ...rest }: OptionProps<SelectOption>) => {
        return (
          <createComponent.Option
            data={data}
            isFocused={isFocused}
            innerRef={innerRef}
            {...rest}
          >
            {data.label || data.value}
            {props?.optionAction &&
              isFocused &&
              props?.optionAction(data.value)}
          </createComponent.Option>
        );
      },
      [props?.optionAction],
    );

    const MultiValueLabel = useCallback(
      (rest: MultiValueGenericProps<SelectOption>) => {
        if (props?.withTooltip) {
          return (
            <createComponent.MultiValueLabel {...rest}>
              <Tooltip
                label={rest.data.label.length > 0 ? rest.data.value : ''}
                side='top'
              >
                {rest.data.label || rest.data.value}
              </Tooltip>
            </createComponent.MultiValueLabel>
          );
        }

        return (
          <createComponent.MultiValueLabel {...rest}>
            {rest.data.label || rest.data.value}
          </createComponent.MultiValueLabel>
        );
      },
      [],
    );

    const MenuList = useCallback((rest: MenuListProps) => {
      return (
        <createComponent.MenuList {...rest}>
          {rest.children}
        </createComponent.MenuList>
      );
    }, []);

    const components = useMemo(
      () => ({
        Control,
        MultiValueLabel,
        MenuList,
        Option: (props?.Option || Option) as ComponentType<OptionProps>,
        DropdownIndicator: () => null,
        ..._components,
      }),
      [Control, MultiValueLabel, _components],
    );
    const defaultClassNames = useMemo(
      () => merge(getDefaultClassNames({ size }), classNames),
      [size, classNames],
    );

    return (
      <AsyncCreatableSelect
        loadOptions={props?.loadOptions}
        cacheOptions
        ref={ref}
        components={components}
        closeMenuOnSelect={false}
        isMulti
        unstyled
        isClearable={false}
        tabSelectsValue={true}
        classNames={defaultClassNames}
        {...props}
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
  multiValue: () => getMultiValueClassNames(''),
  multiValueLabel: () => getMultiValueLabelClassNames('', size),
  multiValueRemove: () => getMultiValueRemoveClassNames('', size),
  groupHeading: () => 'text-gray-400 text-sm px-3 py-1.5 font-normal uppercase',
  valueContainer: () => 'gap-1 py-0.5 mr-0.5',
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
      'flex items-center cursor-pointer overflow-visible',
      props?.isReadOnly && 'pointer-events-none',
      props?.isFocused && 'border-primary-500',
    ),
  });

  return twMerge(defaultStyle, className, variant);
};
