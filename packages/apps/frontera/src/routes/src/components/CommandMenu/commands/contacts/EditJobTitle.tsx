import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { EditJobRoleUseCase } from '@domain/usecases/command-menu/edit-jobTitle.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditJobTitle = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedIds = context.ids;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const handleChangeJobTitle = () => {
    if (!contact) return;

    if (selectedIds?.length === 1) {
      jobRoleUseCase.submitJobRole(
        contact.id,
        contact.value.primaryOrganizationId || '',
      );
    }
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };
  const jobRoleUseCase = useMemo(
    () => new EditJobRoleUseCase(String(contact?.id)),
    [contact?.id],
  );

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        placeholder='Edit job title'
        value={jobRoleUseCase.getJobRole || ''}
        onValueChange={(value) => {
          jobRoleUseCase.setJobRole(value);
        }}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeJobTitle}
        >{`Rename job title to "${
          jobRoleUseCase.getJobRole || ''
        }"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
