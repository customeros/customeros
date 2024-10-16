import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import {
  CreatableSelect,
  CreatableSelectProps,
  getMenuListClassNames,
  getMultiValueClassNames,
} from '@ui/form/CreatableSelect';

interface TagsProps {
  dataTest?: string;
  onBlur?: () => void;
  placeholder: string;
  autofocus?: boolean;
  value: SelectOption[];
  icon?: React.ReactNode;
  closeMenuOnSelect?: boolean;
  menuPortalTarget?: HTMLElement;
  size?: CreatableSelectProps['size'];
  onCreateOption?: (value: string) => void;
  onChange: (value: [SelectOption]) => void;
}

export const Tags = observer(
  ({
    icon,
    size,
    dataTest,
    placeholder,
    onCreateOption,
    value,
    onBlur,
    onChange,
    menuPortalTarget,
    autofocus,
    closeMenuOnSelect,
  }: TagsProps) => {
    const store = useStore();

    const options = store.tags
      ? store.tags
          .toArray()
          .filter((t) => t.value.name !== '')
          .map(
            (tag) =>
              ({
                value: tag.value.id,
                label: tag.value.name,
              } as SelectOption),
          )
      : [];

    return (
      <CreatableSelect
        cacheOptions
        value={value}
        onBlur={onBlur}
        leftElement={icon}
        size={size ?? 'xs'}
        onChange={onChange}
        dataTest={dataTest}
        autoFocus={autofocus}
        backspaceRemovesValue
        defaultOptions={options}
        placeholder={placeholder}
        onCreateOption={onCreateOption}
        menuPortalTarget={menuPortalTarget}
        closeMenuOnSelect={closeMenuOnSelect}
        getNewOptionData={(
          inputValue: string,
          optionLabel: React.ReactNode,
        ) => ({
          label: optionLabel as string,
          value: inputValue,
        })}
        loadOptions={(inputValue: string) =>
          new Promise((resolve) => {
            resolve(
              options.filter((option) =>
                option.label.toLowerCase().includes(inputValue.toLowerCase()),
              ),
            );
          })
        }
        classNames={{
          menuList: () => getMenuListClassNames('w-fit'),
          multiValue: () =>
            getMultiValueClassNames(
              'border-1 border-gray-300 flex items-center rounded-md bg-gray-100 px-0.75 text-gray-500',
            ),
          multiValueRemove: () => 'max-h-4',
          control: () => 'max-h-4',
        }}
      />
    );
  },
);
