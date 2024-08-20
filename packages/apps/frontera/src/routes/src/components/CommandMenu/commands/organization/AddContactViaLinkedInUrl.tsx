import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { getExternalUrl } from '@utils/getExternalLink.ts';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

function validateLinkedInProfileUrl(url: string): boolean {
  const linkedInProfileRegex =
    /^(https:\/\/)?(www\.)?linkedin\.com\/in\/([a-zA-Z0-9\-%]{3,100})\/?$/;

  return linkedInProfileRegex.test(url);
}
export const AddContactViaLinkedInUrl = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [url, setUrl] = useState('');
  const [validationError, setValidationError] = useState(false);

  const entity = store.organizations.value.get((context.ids as string[])?.[0]);
  const label = `Organization - ${entity?.value?.name}`;

  const handleConfirm = () => {
    setValidationError(false);

    const isValidUrl = validateLinkedInProfileUrl(url);

    if (isValidUrl) {
      const formattedUrl = getExternalUrl(url);

      store.contacts.createWithSocial({
        socialUrl: formattedUrl,
        organizationId: (context.ids as string[])?.[0],
      });

      setUrl('');

      store.ui.commandMenu.setOpen(false);
      store.ui.commandMenu.setType('OrganizationCommands');

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
        <CommandItem
          onSelect={handleConfirm}
        >{`Add contact via LinkedIn "${url}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
