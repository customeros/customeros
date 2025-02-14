import { useMemo } from 'react';

import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { EditOrganizationTagUsecase } from '@domain/usecases/command-menu/edit-organization-tag.usecase.ts';

import { Plus } from '@ui/media/icons/Plus.tsx';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeTags = observer(() => {
  const store = useStore();
  const editTags = useMemo(() => new EditOrganizationTagUsecase(), []);

  useModKey('Enter', () => {
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <Command shouldFilter={false} label='Change or add tags...'>
      <CommandInput
        label={editTags.inputLabel}
        value={editTags.searchTerm}
        placeholder='Change or add tags...'
        onValueChange={editTags.setSearchTerm}
        onKeyDownCapture={(e) => {
          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          }
        }}
      />
      <CommandGroup>
        <Command.List>
          {editTags.tagList?.map((tag) => (
            <CommandItem
              key={tag?.id}
              onSelect={() => editTags.select(tag.id)}
              rightAccessory={
                editTags.organizationTags.has(tag.value.name) ? <Check /> : null
              }
              onKeyDown={(e) => {
                if (e.metaKey && e.key === 'Enter') {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.setOpen(false);

                  return;
                }
              }}
            >
              {tag?.value?.name}
            </CommandItem>
          ))}
          {editTags.searchTerm && (
            <CommandItem leftAccessory={<Plus />} onSelect={editTags.create}>
              <span className='text-gray-700 ml-1'>Create new tag:</span>
              <span className='text-gray-500 ml-1'>{editTags.searchTerm}</span>
            </CommandItem>
          )}
        </Command.List>
      </CommandGroup>
    </Command>
  );
});
