import { observer } from 'mobx-react-lite';
import { OrganizationDatum } from '@store/Organizations/Organization.dto';
import { EditLatestOrganizationActive } from '@domain/usecases/command-menu/edit-latest-org.usecase';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
const usecase = new EditLatestOrganizationActive();

export const EditLatestOrgActive = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  const organizations = store.organizations.toArray();
  const contactStore = store.contacts.value.get(context.ids?.[0] as string);

  if (!contactStore) return null;

  const handleChangeOrganization = (value: OrganizationDatum) => {
    contactStore?.draft();

    if (contactStore) {
      contactStore.value.primaryOrganizationId = value.id;
    }

    contactStore?.commit();

    if (contactStore) {
      contactStore?.draft();
      contactStore.value.primaryOrganizationName = value.name;
      contactStore?.commit({ syncOnly: true });
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
        placeholder='Edit organization'
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
