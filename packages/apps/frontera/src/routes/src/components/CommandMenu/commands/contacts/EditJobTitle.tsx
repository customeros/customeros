import { set } from 'lodash';
import { CommandList } from 'cmdk';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const EditJobTitle = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const contactJobTitle =
    contact?.value?.jobRoles?.map((jobRole) => jobRole.jobTitle)?.[0] ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeJobTitle = (jobTitle: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;
    contact?.update((value) => {
      set(value, 'jobRoles[0].jobTitle', jobTitle);

      return value;
    });
  };

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        value={contactJobTitle}
        placeholder='Edit job title'
        onValueChange={handleChangeJobTitle}
      />
      <CommandList></CommandList>
    </Command>
  );
});
