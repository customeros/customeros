import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { SetEmailCase } from '@domain/Contacts/SetEmail.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const useCase = new SetEmailCase();

export const EditEmail = observer(() => {
  const store = useStore();
  const [error, setError] = useState<string | null>('');
  const context = store.ui.commandMenu.context;
  const selectedId = store.ui.selectionId;

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const emailAdress =
    selectedId !== null
      ? contact?.value?.emails?.[selectedId ?? 0]?.email ?? ''
      : contact?.value?.primaryEmail?.email ?? '';

  const [value, setValue] = useState(() => emailAdress);

  const label = `Contact - ${contact?.name}`;

  useEffect(() => {
    if (contact) {
      useCase.setEntity(contact);
    }
  }, [contact?.id]);

  if (!contact) return;

  const primaryEmail = contact.value.primaryEmail?.email;

  const handleSaveEmail = () => {
    if (store.ui.focusRow && !primaryEmail) {
      useCase.setPrimaryEmailForContact(true);
      useCase.setEmailForContact();
      store.ui.setSelectionId(null);
    }

    if (selectedId && store.ui.focusRow) {
      useCase.updateEmailForContact(selectedId);
      store.ui.setSelectionId(null);
    }

    if (selectedId !== null && !store.ui.focusRow) {
      useCase.setPreviousEmail(selectedId);
      useCase.updateEmailForContact(selectedId);
      store.ui.setSelectionId(null);
    }

    if (
      selectedId === null &&
      contact.value.emails.length > 0 &&
      primaryEmail
    ) {
      useCase.updatePrimaryEmailForContact();
      store.ui.setSelectionId(null);
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

  useEffect(() => {
    if (!contact) return;

    if (error) {
      contact.value.emails = contact.value.emails.filter(
        (email) => email.email !== value,
      );
    }
  }, [store.ui.commandMenu.isOpen]);

  return (
    <Command shouldFilter={false}>
      <CommandInput
        label={label}
        value={value}
        placeholder={emailAdress.length > 0 ? 'Edit email' : 'Add new email'}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
        onValueChange={(newValue) => {
          setValue(newValue);
          useCase.setEmail(newValue);

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
          {value ? `Rename email to "${value}"` : `Add new email "${value}"`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
