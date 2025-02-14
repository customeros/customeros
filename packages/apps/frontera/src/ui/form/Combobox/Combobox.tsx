import { cn } from '@ui/utils/cn';
import {
  Select,
  SelectProps,
  getMenuClassNames,
  getOptionClassNames,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select/Select';

interface ComboboxProps extends Omit<SelectProps, 'size'> {
  maxHeight?: string;
}

export const Combobox = ({
  isReadOnly,
  maxHeight,
  isSearchable = true,
  ...props
}: ComboboxProps) => {
  return (
    <Select
      autoFocus
      size='sm'
      menuIsOpen
      backspaceRemovesValue
      isReadOnly={isReadOnly}
      controlShouldRenderValue={false}
      styles={{ menuList: (base) => ({ ...base, maxHeight }) }}
      classNames={{
        input: () => 'pl-3',
        placeholder: () => 'pl-3 text-gray-400',
        container: ({ isFocused }) =>
          getContainerClassNames(
            cn('flex flex-col', isSearchable ? 'pt-2' : 'pt-0 mt-0'),
            'unstyled',
            {
              isFocused,
              size: 'sm',
            },
          ),
        option: ({ isFocused }) =>
          getOptionClassNames('!cursor-pointer', { isFocused }),
        menuList: () =>
          getMenuListClassNames(
            cn('p-0 border-none bg-transparent shadow-none'),
          ),
        menu: ({ menuPlacement }) =>
          getMenuClassNames(menuPlacement)(
            cn('!relative', !isSearchable && 'mt-0'),
            'sm',
          ),
        noOptionsMessage: () => 'text-gray-500',
      }}
      {...props}
      isSearchable={isSearchable}
      onKeyDown={(e) => {
        e.stopPropagation();
        props?.onKeyDown?.(e);
      }}
    />
  );
};
