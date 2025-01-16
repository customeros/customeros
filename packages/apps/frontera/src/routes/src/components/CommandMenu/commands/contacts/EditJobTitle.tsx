import { observer } from 'mobx-react-lite';
import { EditJobRole } from '@domain/usecases/command-menu/edit-jobTitle.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const jobRoleUseCase = new EditJobRole();

export const EditJobTitle = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedIds = context.ids;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const jobRoles = store.contacts.getById(String(contact?.id))?.jobRoles;
  const findPrimaryJobRole = jobRoles?.find(
    (j) => j.primary && j.contact?.metadata.id === contact?.id,
  );

  const jobRolesStore = store.jobRoles.getById(findPrimaryJobRole?.id || '');

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const handleChangeJobTitle = () => {
    if (!contact) return;

    if (selectedIds?.length === 1) {
      jobRoleUseCase.submitJobRole(
        String(contact.id),
        contact.value.primaryOrganizationId || '',
      );
    }
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        placeholder='Edit job title'
        value={findPrimaryJobRole?.jobTitle || ''}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
        onValueChange={(value) => {
          const newValue = value;

          jobRoleUseCase.setJobRole(newValue);

          if (jobRolesStore) {
            jobRolesStore.value.jobTitle = newValue;
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeJobTitle}
        >{`Rename job title to "${
          findPrimaryJobRole?.jobTitle || ''
        }"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
