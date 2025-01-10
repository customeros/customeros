import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';
import { EditPersonaTagUsecase } from '@domain/usecases/command-menu/edit-persona-tag.usecase';

import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const usecase = new EditPersonaTagUsecase();

export const EditPersonaTag = observer(() => {
  const store = useStore();

  const handleSelect = (t: TagDatum) => () => {
    usecase.select(t);
  };

  useModKey('Enter', (e) => {
    e.stopPropagation();
    usecase.reset();
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <Command shouldFilter={false} label='Change or add tags...'>
      <CommandInput
        label={usecase.inputLabel}
        value={usecase.searchTerm}
        placeholder='Edit persona tag...'
        onValueChange={usecase.setSearchTerm}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }

          if (e.metaKey && e.key === 'Enter') {
            e.stopPropagation();
            usecase.reset();
            store.ui.commandMenu.setOpen(false);
          } else {
            // handleSelect(search as unknown as Tag);
          }
        }}
      />

      <CommandGroup>
        <Command.List>
          {usecase.tagList?.map((tag) => (
            <CommandItem
              key={tag.id}
              onSelect={handleSelect(tag.value)}
              rightAccessory={
                usecase.contactTags.has(tag.value.name) ? <Check /> : null
              }
            >
              {tag.value.name}
            </CommandItem>
          ))}
          {usecase.searchTerm && (
            <CommandItem leftAccessory={<Plus />} onSelect={usecase.create}>
              <span className='text-gray-700 ml-1'>Create new tag:</span>
              <span className='text-gray-500 ml-1'>{usecase.searchTerm}</span>
            </CommandItem>
          )}
        </Command.List>
      </CommandGroup>
    </Command>
  );
});
