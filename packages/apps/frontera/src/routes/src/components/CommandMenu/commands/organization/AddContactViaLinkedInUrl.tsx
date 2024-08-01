import React, { useState, MouseEvent } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
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

  const entity = store.organizations.value.get(context.id as string);
  const label = `Organization - ${entity?.value?.name}`;

  const handleConfirm = (e: MouseEvent<HTMLButtonElement> | KeyboardEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setValidationError(false);

    const isValidUrl = validateLinkedInProfileUrl(url);

    if (isValidUrl) {
      const formattedUrl = getExternalUrl(url);

      store.contacts.createWithSocial({
        socialUrl: formattedUrl,
        organizationId: context.id as string,
      });

      setUrl('');

      store.ui.commandMenu.toggle('AddContactViaLinkedInUrl');

      return;
    }
    setValidationError(true);
  };

  useKeyBindings({
    Enter: handleConfirm,
  });

  return (
    <Command label={`Rename `}>
      <div className='p-6 pb-4 flex flex-col gap-2 border-b border-b-gray-100'>
        {label && (
          <Tag size='lg' variant='subtle' colorScheme='gray'>
            <TagLabel>{label}</TagLabel>
          </Tag>
        )}
        <Input
          autoFocus
          size='sm'
          value={url}
          autoComplete='off'
          name='linkedin-input'
          placeholder='Contact`s LinkedIn URL'
          className={cn(validationError && 'border-error-600')}
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
