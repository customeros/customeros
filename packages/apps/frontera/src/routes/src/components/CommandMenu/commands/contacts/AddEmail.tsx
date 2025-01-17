import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddEmailCase } from '@domain/usecases/command-menu/add-email.usecase';

import { Mail02 } from '@ui/media/icons/Mail02';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const addNewEmailCase = new AddEmailCase();

export const AddEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  useEffect(() => {
    if (contact) {
      addNewEmailCase.setEntity(contact);
    }
  }, [contact?.id]);

  const error = addNewEmailCase.error;

  const handleSubmit = async () => {
    await addNewEmailCase.submit();

    const success = !addNewEmailCase.error;

    if (success) {
      store.ui.commandMenu.setOpen(false);
      store.ui.commandMenu.setType('ContactCommands');
    }
  };

  return (
    <Command shouldFilter={false}>
      <CommandInput
        label={label}
        placeholder='Add new email'
        value={addNewEmailCase.inputValue}
        onKeyDownCapture={(e) => {
          if (e.key === '') {
            e.stopPropagation();
          }
        }}
        onValueChange={(value) => {
          addNewEmailCase.setInputValue(value);
          addNewEmailCase.setErrors('');
        }}
      />
      <div className='flex flex-col'>
        {error && <p className='text-error-500 text-sm pl-6'>{error}</p>}{' '}
        <Command.List>
          <CommandItem onSelect={handleSubmit} leftAccessory={<Mail02 />}>
            {`Add new email ${addNewEmailCase.inputValue}`}
          </CommandItem>
        </Command.List>
      </div>
    </Command>
  );
});
