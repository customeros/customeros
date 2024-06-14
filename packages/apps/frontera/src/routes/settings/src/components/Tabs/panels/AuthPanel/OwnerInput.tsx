import { observer } from 'mobx-react-lite';

import { User } from '@graphql/types';
import { Select } from '@ui/form/Select';
import { User02 } from '@ui/media/icons/User02';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
  onSelect: (owner: SelectOption) => void;
}

export const OwnerInput = observer(({ id, onSelect }: OwnerProps) => {
  const store = useStore();
  const users = store.users.toComputedArray((arr) => {
    return arr.filter(
      (e) =>
        Boolean(e.value.firstName) ||
        Boolean(e.value.lastName) ||
        Boolean(e.value.name),
    );
  });

  const options = users
    ?.map((user) => ({
      value: user.id,
      label: user.name,
    }))
    ?.sort((a, b) => a.label.localeCompare(b.label));

  const value = id ? options?.find((o) => o.value === id) : null;

  return (
    <Select
      isClearable
      value={value}
      isLoading={false}
      placeholder='Owner'
      backspaceRemovesValue
      onChange={onSelect}
      options={options}
      leftElement={<User02 className='text-gray-500 mr-3' />}
    />
  );
});
