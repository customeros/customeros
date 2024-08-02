import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { useDidMount, useKeyBindings } from 'rooks';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { Tag, TagLabel } from '@ui/presentation/Tag';

export const RenameOrganizationProperty = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [allowSubmit, setAllowSubmit] = useState(false);
  const entity = store.organizations.value.get(context.ids?.[0] as string);
  const label = `Organization - ${entity?.value?.name}`;
  const property = context.property as 'name' | 'website';

  const handleSelect = (newValue: string) => {
    if (!context.ids?.[0]) return;
    const property = context.property as 'name' | 'website';

    if (!entity || !property) return;

    entity?.update((value) => {
      value[property] = newValue;

      return value;
    });
  };

  useDidMount(() => {
    setTimeout(() => {
      setAllowSubmit(true);
    }, 100);
  });

  const defaultValue = match({ property })
    .with({ property: 'name' }, () => entity?.value?.name ?? '')
    .with({ property: 'website' }, () => entity?.value?.website ?? '')
    .otherwise(() => '');

  const placeholder = match({ property })
    .with({ property: 'name' }, () => 'Rename organization...')
    .with({ property: 'website' }, () => 'Edit website...')
    .otherwise(() => '');

  useKeyBindings(
    {
      Enter: () => {
        store.ui.commandMenu.toggle('RenameOrganizationProperty');
      },
    },
    {
      when: allowSubmit,
    },
  );

  return (
    <Command label={`Rename ${context.property}`}>
      <div className='p-6 pb-4 flex flex-col gap-2 border-b border-b-gray-100'>
        {label && (
          <Tag size='md' variant='subtle' colorScheme='gray'>
            <TagLabel>{label}</TagLabel>
          </Tag>
        )}
        <Input
          autoFocus
          variant='unstyled'
          placeholder={placeholder}
          defaultValue={defaultValue}
          onChange={(e) => handleSelect(e.target.value)}
        />
      </div>
    </Command>
  );
});
