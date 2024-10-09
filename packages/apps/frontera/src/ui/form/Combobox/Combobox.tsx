import { cn } from '@ui/utils/cn';
import {
  Select,
  SelectProps,
  getMenuClassNames,
  getOptionClassNames,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select/Select';

interface ComboboxProps extends SelectProps {
  maxHeight?: string;
}

export const Combobox = ({
  size,
  isReadOnly,
  maxHeight,
  ...props
}: ComboboxProps) => {
  return (
    <Select
      autoFocus
      menuIsOpen
      size={size}
      backspaceRemovesValue
      isReadOnly={isReadOnly}
      controlShouldRenderValue={false}
      classNames={{
        input: () => 'pl-3',
        placeholder: () => 'pl-3 text-gray-400',
        container: ({ isFocused }) =>
          getContainerClassNames('flex flex-col pt-2', 'unstyled', {
            isFocused,
            size,
          }),
        option: ({ isFocused }) =>
          getOptionClassNames('!cursor-pointer', { isFocused }),
        menuList: () =>
          getMenuListClassNames(
            cn('p-0 border-none bg-transparent shadow-none !max-h-[600px]'),
          ),
        menu: ({ menuPlacement }) =>
          getMenuClassNames(menuPlacement)('!relative', size),
        noOptionsMessage: () => 'text-gray-500 p-1',
        valueContainer: () => '!cursor-text',
      }}
      {...props}
    />
  );
};
