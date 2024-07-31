import {
  Select,
  SelectProps,
  getMenuClassNames,
  getOptionClassNames,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select/Select';

interface ComboboxProps extends SelectProps {}

export const Combobox = ({ size, isReadOnly, ...props }: ComboboxProps) => {
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
        option: ({ isFocused, isSelected }) =>
          getOptionClassNames('', { isFocused, isSelected }),
        menuList: () =>
          getMenuListClassNames('p-0 border-none bg-transparent shadow-none'),
        menu: ({ menuPlacement }) =>
          getMenuClassNames(menuPlacement)('!relative', size),
      }}
      {...props}
    />
  );
};
