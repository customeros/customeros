import { useMemo, useState, useEffect } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const [error, setError] = useState<string | null>('');
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

  const [value, setValue] = useState(() => emailAdress);

  const label = `Contact - ${contact?.name}`;

  const handleSaveEmail = () => {
    if (selectedId !== null && !store.ui.focusRow) {
      contact?.updateEmail({
        previousEmail: oldEmail ?? '',
        index: selectedId,
        primary: false,
      });
      store.ui.commandMenu.setOpen(false);
      store.ui.setSelectionId(null);
      store.ui.commandMenu.setType('ContactCommands');
    }

    if (selectedId === null) {
      contact?.updateEmailPrimary(oldEmail ?? '');
      store.ui.commandMenu.setOpen(false);
      store.ui.setSelectionId(null);
      store.ui.commandMenu.setType('ContactCommands');
    }

    if (store.ui.focusRow) {
      contact?.updateEmail(
        {
          previousEmail: '',
          index: selectedId ?? 0,
          primary: false,
        },
        {
          onError: (error) => {
            setError(error);
          },
          invalidate: false,
          onSuccess(serverId) {
            if (serverId) {
              store.ui.commandMenu.setOpen(false);
              store.ui.setSelectionId(null);
              store.ui.commandMenu.setType('ContactCommands');
            }
          },
        },
      );
    }

    if (
      store.ui.focusRow &&
      !contact?.value.primaryEmail?.email &&
      selectedId === null
    ) {
      contact?.updateEmailPrimary('');
      store.ui.commandMenu.setOpen(false);
      store.ui.setSelectionId(null);
      store.ui.commandMenu.setType('ContactCommands');
    }
  };

  useEffect(() => {
    if (store.ui.commandMenu.isOpen === false) {
      store.ui.setSelectionId(null);
    }
  }, [store.ui.commandMenu.isOpen]);

  useEffect(() => {
    if (!contact) return;

    if (error) {
      contact.value.emails = contact.value.emails.filter(
        (email) => email.email !== value,
      );
    }
  }, [store.ui.commandMenu.isOpen]);

  return (
    <Command>
      <CommandInput
        label={label}
        value={emailAdress ?? value}
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
                setValue(newValue);
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

          if (error) {
            setError('');
          }
        }}
      />
      {error && (
        <p className='ml-5 text-xs text-error-600 mt-2'>
          This email is already used by another contact
        </p>
      )}
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
