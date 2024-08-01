import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { Tag, TagLabel } from '@ui/presentation/Tag';

export const RenameOrganizationProperty = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = store.organizations.value.get(context.id as string);
  const label = `Organization - ${entity?.value?.name}`;
  const property = context.property as 'name' | 'website';

  const handleSelect = (newValue: string) => {
    if (!context.id) return;
    const property = context.property as 'name' | 'website';

    if (!entity || !property) return;

    entity?.update((value) => {
      value[property] = newValue;
      return value;
    });
  };

  const defaultValue = match({ property })
    .with({ property: 'name' }, () => entity?.value?.name ?? '')
    .with({ property: 'website' }, () => entity?.value?.website ?? '')
    .otherwise(() => '');

  const placeholder = match({ property })
    .with({ property: 'name' }, () => 'Rename organization...')
    .with({ property: 'website' }, () => 'Edit website...')
    .otherwise(() => '');

  useKeyBindings({
    Enter: () => {
      store.ui.commandMenu.clearContext();
      store.ui.commandMenu.toggle('RenameOrganizationProperty');
    },
  });

  return (
    <Command label={`Rename `}>
      <div className='p-6 pb-4 flex flex-col gap-2 border-b border-b-gray-100'>
        {label && (
          <Tag size='lg' variant='subtle' colorScheme='gray'>
            <TagLabel>{label}</TagLabel>
          </Tag>
        )}
        <Input
          autoFocus
          placeholder={placeholder}
          defaultValue={defaultValue}
          onChange={(e) => handleSelect(e.target.value)}
        />
      </div>
    </Command>
  );
});
