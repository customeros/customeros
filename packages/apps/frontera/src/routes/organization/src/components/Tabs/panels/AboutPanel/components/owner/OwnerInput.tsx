import { observer } from 'mobx-react-lite';

import { User } from '@graphql/types';
import { Select } from '@ui/form/Select';
import { Key01 } from '@ui/media/icons/Key01';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';

type Owner = Pick<User, 'id' | 'firstName' | 'lastName'> | null;
interface OwnerProps {
  id: string;
  owner?: Owner;
  dataTest?: string;
}

export const OwnerInput = observer(({ id, owner, dataTest }: OwnerProps) => {
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

  const value = owner ? options?.find((o) => o.value === owner.id) : null;

  const handleSelect = (option: SelectOption) => {
    const organization = store.organizations.value.get(id);
    const targetOwner = store.users.value.get(option?.value);

    organization?.update((value) => {
      if (!option || !option?.value) {
        value.owner = null;
      } else {
        value.owner = targetOwner?.value;
      }

      return value;
    });
  };

  return (
    <Select
      isClearable
      value={value}
      isLoading={false}
      options={options}
      placeholder='Owner'
      dataTest={dataTest}
      backspaceRemovesValue
      onChange={handleSelect}
      leftElement={<Key01 className='text-gray-500 mr-3' />}
    />
  );
});
