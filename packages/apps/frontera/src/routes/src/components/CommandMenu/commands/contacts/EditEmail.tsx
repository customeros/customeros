import { useMemo, useEffect } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedId = store.ui.selectionId;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const oldEmail = useMemo(
    () =>
      contact?.value?.emails?.[selectedId ?? 0]?.email ||
      contact?.value?.primaryEmail?.email,
    [contact?.isLoading, selectedId],
  );
  const emailAdress =
    selectedId !== null
      ? contact?.value?.emails?.[selectedId ?? 0]?.email ?? ''
      : contact?.value?.primaryEmail?.email ?? '';

  const label = `Contact - ${contact?.name}`;

  const handleSaveEmail = () => {
    if (selectedId !== null && !store.ui.focusRow) {
      contact?.updateEmail(oldEmail ?? '', selectedId ?? 0);
    }

    if (selectedId === null) {
      contact?.updateEmailPrimary(oldEmail ?? '');
    }

    if (store.ui.focusRow) {
      contact?.updateEmail('', selectedId ?? 0);
    }

    if (
      store.ui.focusRow &&
      !contact?.value.primaryEmail?.email &&
      selectedId === null
    ) {
      contact?.updateEmailPrimary('');
    }
    store.ui.commandMenu.setOpen(false);
    store.ui.setSelectionId(null);
    store.ui.commandMenu.setType('ContactCommands');
  };

  useEffect(() => {
    if (store.ui.commandMenu.isOpen === false) {
      store.ui.setSelectionId(null);
    }
  }, [store.ui.commandMenu.isOpen]);

  return (
    <Command>
      <CommandInput
        label={label}
        value={emailAdress}
        placeholder={emailAdress.length > 0 ? 'Edit email' : 'Add new email'}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
        onValueChange={(newValue) => {
          contact?.update(
            (value) => {
              if (selectedId !== null) {
                set(value, ['emails', selectedId ?? 0, 'email'], newValue);
              } else {
                if (newValue.length === 0) {
                  set(value, 'primaryEmail', null);
                } else {
                  set(value, 'primaryEmail.email', newValue);
                }
              }

              return value;
            },
            { mutate: false },
          );
        }}
      />
      <Command.List>
        <CommandItem leftAccessory={<Edit03 />} onSelect={handleSaveEmail}>
          {(oldEmail ?? '').length > 0
            ? `Rename email to "${emailAdress}"`
            : `Add new email "${emailAdress}"`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
