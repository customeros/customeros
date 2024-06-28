import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import {
  CreatableSelect,
  getMenuListClassNames,
  getContainerClassNames,
  getMultiValueClassNames,
} from '@ui/form/CreatableSelect';

interface TagsProps {
  placeholder: string;
  autofocus?: boolean;
  hideBorder?: boolean;
  icon: React.ReactNode;
  value: SelectOption[];
  menuPortalTarget?: HTMLElement;
  onCreateOption?: (value: string) => void;
  onChange: (value: [SelectOption]) => void;
}

export const Tags = observer(
  ({
    icon,
    placeholder,
    onCreateOption,
    value,
    onChange,
    menuPortalTarget,
    autofocus,
    hideBorder,
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
        autoFocus={autofocus}
        onChange={onChange}
        backspaceRemovesValue
        menuPortalTarget={menuPortalTarget}
        defaultOptions={options}
        placeholder={placeholder}
        onCreateOption={onCreateOption}
        leftElement={icon}
        classNames={{
          menuList: () => getMenuListClassNames('w-fit'),
          multiValue: () =>
            getMultiValueClassNames(
              'border-1 border-gray-300 rounded-full bg-gray-100 text-gray-500',
            ),
          container: () =>
            hideBorder
              ? getContainerClassNames('', 'unstyled')
              : getContainerClassNames(
                  'border-b border-gray-300 focus-within:border-primary-600',
                  'unstyled',
                ),
        }}
        loadOptions={(inputValue: string) =>
          new Promise((resolve) => {
            resolve(
              options.filter((option) =>
                option.label.toLowerCase().includes(inputValue.toLowerCase()),
              ),
            );
          })
        }
        getNewOptionData={(
          inputValue: string,
          optionLabel: React.ReactNode,
        ) => ({
          label: optionLabel as string,
          value: inputValue,
        })}
      />
    );
  },
);
