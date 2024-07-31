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
  const label = `organization ${context.property} - ${entity?.value?.name}`;

  const handleSelect = (newValue: string) => {
    if (!context.id) return;
    const property = context.property as 'name' | 'website' | 'description';

    if (!entity || !property) return;

    entity?.update((value) => {
      value[property] = newValue;
      return value;
    });
  };

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
            <TagLabel>Update {label}</TagLabel>
          </Tag>
        )}
        <Input
          onChange={(e) => handleSelect(e.target.value)}
          placeholder={`Type a new name for ${entity?.value?.name}`}
        />
      </div>
    </Command>
  );
});
