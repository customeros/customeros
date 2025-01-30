import { observer } from 'mobx-react-lite';
import { OrganizationDatum } from '@store/Organizations/Organization.dto';
import { EditLatestOrganizationActive } from '@domain/usecases/command-menu/edit-latest-org.usecase';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
const usecase = new EditLatestOrganizationActive();

export const EditLatestOrgActive = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.getById(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  const organizations = store.organizations.toArray();

  if (!contact) return null;

  const handleChangeOrganization = (value: OrganizationDatum) => {
    if (contact) {
      contact?.draft();
      contact.value.primaryOrganizationName = value.name;
      contact.value.primaryOrganizationId = value.id;
      contact?.commit();
    }

    const orgStore = store.organizations.getById(value.id);

    orgStore?.draft();
    orgStore?.value.contacts.push(contact?.id || '');
    orgStore?.commit({ syncOnly: true });

    handleClose();
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={usecase.searchTerm}
        placeholder='Change company'
        onValueChange={usecase.setSearchTerm}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        {organizations.map((option) => (
          <CommandItem
            key={option.id}
            onSelect={() =>
              handleChangeOrganization(option as unknown as OrganizationDatum)
            }
          >
            {option.name}
          </CommandItem>
        ))}
      </Command.List>
    </Command>
  );
});
