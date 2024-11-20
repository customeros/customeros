import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getExternalUrl } from '@utils/getExternalLink';
import { validLinkedInProfileUrl } from '@utils/linkedinValidation';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const AddContactViaLinkedInUrl = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [url, setUrl] = useState('');
  const [validationError, setValidationError] = useState(false);

  const entity = store.organizations.value.get((context.ids as string[])?.[0]);

  const label = match(context.entity)
    .with('Organization', () => `Organization - ${entity?.value?.name}`)
    .with('Contact', () => 'Contact')
    .otherwise(() => '');

  const handleConfirm = () => {
    setValidationError(false);

    const isValidUrl = validLinkedInProfileUrl(url);

    if (isValidUrl) {
      const formattedUrl = getExternalUrl(url);

      match(context.entity)
        .with('Organization', () => {
          store.contacts.createWithSocial({
            socialUrl: formattedUrl,
            organizationId: (context.ids as string[])?.[0],
          });
          store.ui.commandMenu.setOpen(false);
          store.ui.commandMenu.setType('OrganizationCommands');
        })
        .with('Contact', () => {
          store.contacts.createWithoutOrg({
            socialUrl: formattedUrl,
          });

          store.ui.commandMenu.setOpen(false);
          store.ui.commandMenu.setType('ContactCommands');
        });

      setUrl('');

      return;
    }
    setValidationError(true);
  };

  return (
    <Command label={`Add contact via LinkedIn`}>
      <CommandInput
        value={url}
        label={label}
        placeholder='Add contact via LinkedIn'
        onValueChange={(value) => setUrl(value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            handleConfirm();
          }
        }}
        onKeyUp={(e) => {
          if (e.key === 'Backspace' && url.length === 0) {
            setValidationError(false);
          }
        }}
      />
      {validationError && (
        <p className='ml-5 text-xs text-error-600 mt-2'>
          Enter a valid LinkedIn profile URL (e.g. linkedin.com/in/identifier)
        </p>
      )}
      <Command.List>
        <CommandItem className='' onSelect={handleConfirm}>
          <span className='overflow-hidden text-ellipsis whitespace-nowrap'>{`Add contact via LinkedIn "${url}"`}</span>
        </CommandItem>
      </Command.List>
    </Command>
  );
});
