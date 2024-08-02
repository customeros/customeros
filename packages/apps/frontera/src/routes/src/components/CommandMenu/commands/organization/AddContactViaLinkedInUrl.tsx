import React, { useState, MouseEvent } from 'react';

import { observer } from 'mobx-react-lite';
import { useDidMount, useKeyBindings } from 'rooks';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Command } from '@ui/overlay/CommandMenu';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { getExternalUrl } from '@utils/getExternalLink.ts';

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
  const [allowSubmit, setAllowSubmit] = useState(false);

  const entity = store.organizations.value.get((context.ids as string[])?.[0]);
  const label = `Organization - ${entity?.value?.name}`;

  useDidMount(() => {
    setTimeout(() => {
      setAllowSubmit(true);
    }, 100);
  });

  const handleConfirm = (e: MouseEvent<HTMLButtonElement> | KeyboardEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setValidationError(false);

    const isValidUrl = validateLinkedInProfileUrl(url);

    if (isValidUrl) {
      const formattedUrl = getExternalUrl(url);

      store.contacts.createWithSocial({
        socialUrl: formattedUrl,
        organizationId: (context.ids as string[])?.[0],
      });

      setUrl('');

      store.ui.commandMenu.toggle('AddContactViaLinkedInUrl');

      return;
    }
    setValidationError(true);
  };

  useKeyBindings(
    {
      Enter: handleConfirm,
    },
    {
      when: allowSubmit,
    },
  );

  return (
    <Command label={`Rename `}>
      <div className='p-6 pb-4 flex flex-col gap-2 border-b border-b-gray-100'>
        {label && (
          <Tag size='md' variant='subtle' colorScheme='gray'>
            <TagLabel>{label}</TagLabel>
          </Tag>
        )}
        <Input
          autoFocus
          value={url}
          autoComplete='off'
          variant='unstyled'
          name='linkedin-input'
          placeholder='Add contact'
          onChange={(e) => {
            setUrl(e.target.value);
          }}
          onKeyDown={(e) => {
            if (e.key === '/') {
              e.stopPropagation();
            }
          }}
        />
        {validationError && (
          <p className='text-xs text-error-600'>
            Enter a valid LinkedIn profile URL (e.g. linkedin.com/in/identifier)
          </p>
        )}
      </div>
    </Command>
  );
});
