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
          getContainerClassNames('flex flex-col pt-2', 'unstyled', {
            isFocused,
            size: 'sm',
          }),
        option: ({ isFocused }) =>
          getOptionClassNames('!cursor-pointer', { isFocused }),
        menuList: () =>
          getMenuListClassNames(
            cn('p-0 border-none bg-transparent shadow-none'),
          ),
        menu: ({ menuPlacement }) =>
          getMenuClassNames(menuPlacement)('!relative', 'sm'),
        noOptionsMessage: () => 'text-gray-500',
        valueContainer: () => '!cursor-text',
      }}
      {...props}
    />
  );
};
