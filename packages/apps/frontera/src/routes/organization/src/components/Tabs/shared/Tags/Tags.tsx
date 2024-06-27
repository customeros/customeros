import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import {
  CreatableSelect,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/CreatableSelect';

interface TagsProps {
  placeholder: string;
  icon: React.ReactNode;
  value: SelectOption[];
  onCreateOption?: (value: string) => void;
  onChange: (value: [SelectOption]) => void;
}

export const Tags = observer(
  ({ icon, placeholder, onCreateOption, value, onChange }: TagsProps) => {
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
        onChange={onChange}
        backspaceRemovesValue
        menuPortalTarget={document.body}
        defaultOptions={options}
        placeholder={placeholder}
        onCreateOption={onCreateOption}
        leftElement={icon}
        classNames={{
          menuList: () => getMenuListClassNames('w-fit'),
          container: () => getContainerClassNames('', 'unstyled'),
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
