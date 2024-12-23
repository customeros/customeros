import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink.ts';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditCompanyLinkedin = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedIds = context.ids;
  const organization = store.organizations.getById(context.ids?.[0]);

  const organizationLinkedIn = organization?.value?.socialMedia.find((social) =>
    social.url.includes('linkedin'),
  );
  const link = organizationLinkedIn?.url ?? '';
  const alias = organizationLinkedIn?.alias;
  const formattedLink = getFormattedLink(link).replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '/',
  );

  const displayLink = alias ? `/${alias}` : formattedLink;
  const url = formattedLink
    ? link?.includes('linkedin')
      ? getExternalUrl(`https://linkedin.com/company${displayLink}`)
      : getExternalUrl(link)
    : '';

  const [name, setName] = useState(() => url);

  const label =
    selectedIds?.length === 1
      ? `Organization - ${organization?.value?.name}`
      : `${selectedIds?.length} contacts`;

  const handleAddSocial = () => {
    const url = name;

    if (!url) return;

    if (!organization || url === 'Unknown' || url === '') return;

    const formattedValue =
      url.includes('https://www') || url.includes('linkedin.com')
        ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
        : `in/${url}`;

    organization.draft();
    organization.addSocial(`linkedin.com/${formattedValue}`);
    organization.commit();
  };

  const handleUpdateSocial = () => {
    const url = name;

    if (!url) return;

    if (!organization?.value) return;
    const linkedinId = organization?.value?.socialMedia.find((social) =>
      social.url.includes('linkedin'),
    )?.id;

    if (!linkedinId) return;

    const idx = organization.value?.socialMedia.findIndex(
      (s) => s.id === linkedinId,
    );

    organization.draft();

    if (idx !== -1) {
      const formattedValue =
        url.includes('https://www') || url.includes('linkedin.com')
          ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
          : `in/${url}`;

      organization.value.socialMedia[
        idx
      ].url = `linkedin.com/${formattedValue}`;
    }

    if (url === '') {
      organization.value.socialMedia.splice(idx, 1);
    }

    organization.commit();
  };

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        value={name || ''}
        placeholder='Edit LinkedIn'
        onValueChange={(value) => setName(value)}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={
            organizationLinkedIn
              ? () => handleUpdateSocial()
              : () => handleAddSocial()
          }
        >{`Change to: "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
