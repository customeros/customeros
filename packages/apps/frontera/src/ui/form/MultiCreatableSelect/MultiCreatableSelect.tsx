import { useMemo, forwardRef, useCallback, ComponentType } from 'react';
import {
  OptionProps,
  ControlProps,
  MenuListProps,
  SelectInstance,
  MultiValueProps,
  ClassNamesConfig,
  MultiValueGenericProps,
  components as createComponent,
} from 'react-select';

import merge from 'lodash/merge';
import { twMerge } from 'tailwind-merge';
import AsyncCreatableSelect, {
  AsyncCreatableProps,
} from 'react-select/async-creatable';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

import { SelectProps } from '../Select';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
  withTooltip?: boolean;
  size?: 'xs' | 'sm' | 'md' | 'lg';
  navigateAfterAddingToContract?: boolean;
  removeValue?: (value: SelectOption) => void;
  optionAction?: (data: string) => JSX.Element;
  Option?: ComponentType<OptionProps<SelectOption>>;
  MultiValue?: ComponentType<MultiValueProps<SelectOption>>;
}

export const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  (
    {
      size = 'md',
      name,
      formId,
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
            className={`flex w-full items-center group ${sizeClass}`}
            {...innerProps}
          >
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
  isReadOnly,
}: Pick<SelectProps, 'size' | 'isReadOnly'>): ClassNamesConfig => ({
  container: ({ isFocused }) =>
    twMerge(
      'flex mt-1 cursor-pointer overflow-visible min-w-[300px] w-full focus-visible:border-0 focus:border-0',
      isReadOnly && 'pointer-events-none',
      isFocused && 'border-primary-500',
    ),
  menu: ({ menuPlacement }) =>
    cn(
      menuPlacement === 'top'
        ? 'mb-2 animate-slideUpAndFade'
        : 'mt-2 animate-slideDownAndFade',
      'z-50',
    ),
  menuList: () =>
    'p-2 z-50  max-h-[12rem] border border-gray-200 bg-white rounded-lg shadow-lg overflow-y-auto overscroll-auto',
  option: ({ isFocused, isSelected }) =>
    cn(
      'my-[2px] px-3 py-1.5 rounded-md text-gray-700 line-clamp-1 text-sm transition ease-in-out delay-50 hover:bg-gray-50',
      isSelected && 'bg-gray-50 font-medium leading-normal',
      isFocused && 'ring-2 ring-gray-100',
    ),
  placeholder: () => 'text-gray-400 text-inherit',
  multiValue: () =>
    'flex p-0 gap-0 text-gray-700 text-sm mr-1 cursor-default h-[auto]',
  multiValueLabel: () => 'text-gray-700 text-sm mr-1 h-[20px] self-center',
  multiValueRemove: () => 'hidden',
  groupHeading: () => 'text-gray-400 text-sm px-3 py-1.5 font-normal uppercase',
  valueContainer: () => getValueContainerClassNames(),
  control: () => 'overflow-visible',
  input: () => 'overflow-visible text-gray-500 leading-4',
});

export const getMultiValueClassNames = (className?: string) => {
  const defaultStyle =
    'flex p-0 gap-0 text-gray-700 text-sm mr-1 cursor-default h-4';

  return twMerge(defaultStyle, className);
};

export const getMultiValueLabelClassNames = (className?: string) => {
  const defaultStyle = 'text-gray-700 text-sm mr-1 h-[20px] self-center';

  return twMerge(defaultStyle, className);
};

export const getMenuListClassNames = (className?: string) => {
  const defaultStyle =
    'p-2 z-50 max-h-[12rem] border border-gray-200 bg-white rounded-lg shadow-lg overflow-y-auto overscroll-auto';

  return twMerge(defaultStyle, className);
};

export const getValueContainerClassNames = (className?: string) => {
  const defaultStyle =
    'overflow-visible max-h-[86px] flex items-center justify-start';

  return twMerge(defaultStyle, className);
};
