import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

import set from 'lodash/set';
import { useKeyBindings } from 'rooks';
import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Spinner } from '@ui/feedback/Spinner';
import { Star06 } from '@ui/media/icons/Star06';
import { IconButton } from '@ui/form/IconButton';
import { getTimezone } from '@utils/getTimezone';
import { useStore } from '@shared/hooks/useStore';
import { Tags } from '@organization/components/Tabs';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { getFormattedLink } from '@utils/getExternalLink';
import { Tag, Social, TableViewType } from '@graphql/types';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';

import { EmailsSection } from './components';
import { EnrichContactModal } from './components/EnrichContactModal';

export const ContactPreviewCard = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const [isOpen, setIsOpen] = useState(false);
  const [isEditName, setIsEditName] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const contactId = store.ui.focusRow;
  const preset = searchParams?.get('preset');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;

  if (tableType !== TableViewType.Contacts && !contactId) {
    store.ui.setContactPreviewCardOpen(false);

    return null;
  }

  if (!contactId) return;

  const contact = store.contacts
    .toArray()
    .find((c) => c.id === String(contactId));

  const fullName = contact?.name || 'Unnamed';
  const src = contact?.value?.profilePhotoUrl;

  const company =
    contact?.value.latestOrganizationWithJobRole?.organization.name;

  const role = contact?.value.latestOrganizationWithJobRole?.jobRole.jobTitle;

  const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;
  const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
  const flag = flags[countryA2 || ''];
  const city = contact?.value.locations?.[0]?.locality;
  const timezone = city
    ? cityTimezone.lookupViaCity(city).find((c) => {
        return c.iso2 === contact.value.locations?.[0].countryCodeA2;
      })?.timezone
    : null;

  const fromatedUrl = contact?.value?.socials?.[0]?.url.replace(
    'https://www.',
    '',
  );
  const href = fromatedUrl?.startsWith('http')
    ? fromatedUrl
    : `https://${fromatedUrl}`;

  const formatedFollowersCount = contact?.value?.socials?.[0]?.followersCount
    ?.toLocaleString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  const handleUpdateSocial = (url: string) => {
    const linkedinId = contact?.value.socials.find((social) =>
      social.url.includes('linkedin'),
    )?.id;

    if (fromatedUrl === undefined && url.trim() !== '') {
      contact?.update((contactData) => {
        const formattedValue =
          url.includes('https://www') || url.includes('linkedin.com')
            ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
            : url;

        contactData.socials.push({
          id: crypto.randomUUID(),
          url: `linkedin.com/${formattedValue}`,
        } as Social);

        return contactData;
      });
    }

    contact?.update((org) => {
      const idx = org.socials.findIndex((s) => s.id === linkedinId);
      const formattedValue =
        url.includes('https://www') || url.includes('linkedin.com')
          ? getFormattedLink(url).replace(/^linkedin\.com\//, '')
          : `in/${url}`;

      if (idx !== -1) {
        org.socials[idx].url = `linkedin.com/${formattedValue}`;
      }

      if (url === '') {
        org.socials.splice(idx, 1);
      }

      return org;
    });
  };

  useKeyBindings(
    {
      Escape: () => {
        store.ui.setContactPreviewCardOpen(false);
      },
      Space: (e) => {
        e.preventDefault();
        store.ui.setContactPreviewCardOpen(false);
      },
    },
    {
      when: store.ui.contactPreviewCardOpen,
    },
  );

  const userBeenEnriched =
    contact?.value?.enrichDetails.enrichedAt ||
    contact?.value?.enrichDetails.failedAt;

  const requestedEnrichment = contact?.value?.enrichDetails.requestedAt;

  return (
    <>
      {store.ui.contactPreviewCardOpen && (
        <div
          data-state={store.ui.contactPreviewCardOpen ? 'open' : 'closed'}
          className='data-[state=open]:animate-slideLeftAndFade data-[state=closed]:animate-slideRightAndFade flex flex-col absolute right-[12px] -top-[-53px] p-4 max-w-[390px] min-w-[350px] border border-gray-200 rounded-lg z-[1] bg-white'
        >
          <div className='flex justify-between items-start'>
            <Avatar
              size='sm'
              textSize='xs'
              name={fullName}
              variant='circle'
              src={src || undefined}
            />
            <div className='flex items-center gap-2'>
              {userBeenEnriched && (
                <Tooltip asChild label='Enrich this contact'>
                  {requestedEnrichment ? (
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      icon={<Star06 />}
                      onClick={() => setIsOpen(true)}
                      aria-label='enrich this contact'
                    />
                  ) : (
                    <Spinner
                      size='sm'
                      label='enriching'
                      className='text-gray-400 fill-gray-700'
                    />
                  )}
                </Tooltip>
              )}
              <IconButton
                size='xxs'
                icon={<X />}
                variant='ghost'
                aria-label='close'
                onClick={() => store.ui.setContactPreviewCardOpen(false)}
              />
            </div>
          </div>
          {isEditName ? (
            <Input
              value={fullName}
              variant='unstyled'
              placeholder='Unknown'
              className='mb-[-8px]'
              onFocus={(e) => e.target.select()}
              onBlur={() => setIsEditName(false)}
              onChange={(e) => {
                contact?.update((value) => {
                  set(value, 'name', e.target.value);

                  return value;
                });
              }}
            />
          ) : (
            <div className='flex h-fit w-full'>
              <p
                onClick={() => setIsEditName(true)}
                className='font-medium mt-2 w-fit  '
              >
                {fullName}
              </p>
              {company ? (
                <div className='flex flex-2 items-center gap-1 w-full'>
                  <p className='mt-2 text-gray-500 ml-1'>at</p>
                  <p className='font-medium mt-2 line-clamp-1'>
                    {company || 'No org yet'}
                  </p>
                </div>
              ) : (
                <span className='mt-2'>(No org yet)</span>
              )}
            </div>
          )}
          <Input
            size='xs'
            variant='unstyled'
            value={role || ''}
            placeholder='Enter title'
            onFocus={(e) => e.target.select()}
            className='w-[290px] overflow-hidden text-ellipsis whitespace-nowrap'
            onChange={(e) => {
              contact?.update((value) => {
                set(value, 'jobRoles[0].jobTitle', e.target.value);

                return value;
              });
            }}
          />
          <div className={cn('flex items-center mb-4', countryA3 && 'gap-1')}>
            <span className='mb-1'>{flag}</span>
            {countryA3 && <span className='ml-2 text-sm'>{countryA3}</span>}
            {countryA3 && city && timezone && <span>•</span>}
            {city && (
              <span className='overflow-hidden text-ellipsis whitespace-nowrap text-sm'>
                {city}
              </span>
            )}
            {city && timezone && <span>•</span>}
            {timezone && (
              <span className='w-[150px] text-sm'>
                {getTimezone(timezone || '')} local time
              </span>
            )}
          </div>
          <div className='flex justify-between gap-1 w-full mb-4 flex-col'>
            <EmailsSection contactId={contactId} />
          </div>

          <div className='flex justify-between gap-1 w-full mb-4'>
            <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
              <Tag01 className='mt-[1px] text-gray-500' />
              Persona
            </div>
            <Tags
              placeholder='No tags yet'
              value={
                contact?.value?.tags?.map((tag: Tag) => ({
                  value: tag.id,
                  label: tag.name,
                })) || []
              }
              onChange={(e) => {
                contact?.update((c) => {
                  c.tags =
                    (e
                      .map((tag) => store.tags?.value.get(tag.value)?.value)
                      .filter(Boolean) as Array<Tag>) ?? [];

                  return c;
                });
              }}
            />
          </div>
          <div className='flex flex-col gap-4'>
            <div className='flex gap- items-center w-full '>
              <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500 ' />
                LinkedIn
              </div>
              <div
                onMouseEnter={() => setIsHovered(true)}
                onMouseLeave={() => setIsHovered(false)}
                className='flex items-center gap-1 w-full'
              >
                <Input
                  size='xs'
                  variant='unstyled'
                  value={fromatedUrl}
                  className='text-ellipsis'
                  onFocus={(e) => e.target.select()}
                  placeholder='LinkedIn profile link'
                  onChange={(e) => {
                    handleUpdateSocial(e.target.value);
                  }}
                />
                {fromatedUrl && isHovered && (
                  <Link to={href} target='_blank'>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      colorScheme='gray'
                      aria-label='social link'
                      className='hover:bg-gray-200 '
                      icon={<LinkExternal02 className='text-gray-500' />}
                    />
                  </Link>
                )}
              </div>
            </div>
            <div className='flex gap-1 w-full'>
              <div className='flex items-center gap-2 mr-[42px] text-sm text-gray-500 '>
                <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                Followers
              </div>
              <span
                className={cn(
                  'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm',
                  formatedFollowersCount ? 'text-gray-700' : 'text-gray-400',
                )}
              >
                {formatedFollowersCount || 'Unknown'}
              </span>
            </div>
            <div className='flex gap-1 w-full mt-[2px]'>
              <div className='flex items-center gap-2 mr-[19px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                Connected to
              </div>
              <span
                className={cn(
                  'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm',
                  contact?.value?.connectedUsers?.[0]?.name
                    ? 'text-gray-700'
                    : 'text-gray-400',
                )}
              >
                {contact?.value?.connectedUsers?.[0]?.name || 'No one yet'}
              </span>
            </div>
            {contact?.value?.enrichDetails.enrichedAt && (
              <div className='bg-grayModern-50 w-full rounded-[4px] border-[1px] border-grayModern-100 px-2 py-1'>
                <p className='text-sm text-center'>{`Last enriched on ${DateTimeUtils.format(
                  contact?.value.enrichDetails.enrichedAt,
                  DateTimeUtils.dateWithHourWithQomma,
                )} `}</p>
              </div>
            )}
          </div>
        </div>
      )}
      <EnrichContactModal
        isModalOpen={isOpen}
        contactId={contactId}
        onClose={() => {
          setIsOpen(false);
        }}
      />
    </>
  );
});
