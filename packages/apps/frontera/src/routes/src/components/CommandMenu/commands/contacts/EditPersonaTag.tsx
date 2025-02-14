import { useMemo } from 'react';

import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { EditPersonaTagUsecase } from '@domain/usecases/command-menu/edit-persona-tag.usecase';

import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditPersonaTag = observer(() => {
  const usecase = useMemo(() => new EditPersonaTagUsecase(), []);

  return (
    <Command shouldFilter={false} label='Change or add tags...'>
      <CommandInput
        label={usecase.inputLabel}
        value={usecase.searchTerm}
        placeholder='Change or add tags...'
        onValueChange={usecase.setSearchTerm}
        onKeyDownCapture={(e) => {
          if (e.metaKey) {
            usecase.allowClose();
          }
        }}
        onKeyUpCapture={(e) => {
          if (e.metaKey) {
            usecase.preventClose();
          }
        }}
      />

      <CommandGroup>
        <Command.List>
          {usecase.tagList?.map((tag) => (
            <CommandItem
              key={tag.id}
              value={tag.value.metadata.id}
              onSelect={() => usecase.select(tag.id)}
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
