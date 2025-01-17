import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

import { useKeyBindings } from 'rooks';
import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';
import { EditJobRole } from '@domain/usecases/contact-preview-card/edit-jobrole.usecase';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { IconButton } from '@ui/form/IconButton';
import { getTimezone } from '@utils/getTimezone';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TableViewType } from '@graphql/types';
import { Tags } from '@organization/components/Tabs';
import { getFormattedLink } from '@utils/getExternalLink';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

import { EmailsSection } from './components';
import { EnrichContactModal } from './components/EnrichContactModal';

const jobRoleUseCase = new EditJobRole();

export const ContactPreviewCard = observer(() => {
  const store = useStore();
  const [_, copyToClipboard] = useCopyToClipboard();

  const [searchParams] = useSearchParams();
  const [isOpen, setIsOpen] = useState(false);
  const [isEditName, setIsEditName] = useState(false);
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

  const company = contact?.value.primaryOrganizationName;

  const jobRoles = store.contacts.getById(String(contactId))?.jobRoles;

  const findPrimaryJobRole = jobRoles?.find(
    (j) => j.primary && j.contact?.metadata.id === contactId,
  );

  const jobRolesStore = store.jobRoles.getById(findPrimaryJobRole?.id || '');

  const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;
  const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
  const flag = flags[countryA2 || ''];
  const city = contact?.value.locations?.[0]?.locality;
  const timezone = city
    ? cityTimezone.lookupViaCity(city).find((c) => {
        return c.iso2 === contact.value.locations?.[0].countryCodeA2;
      })?.timezone
    : null;

  const linkedInProfile = contact?.value.linkedInUrl;
  const fromatedUrl = getFormattedLink(linkedInProfile || '').replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '/',
  );
  const href = contact?.value.linkedInUrl;

  const formatedFollowersCount = contact?.value?.linkedInFollowerCount
    ?.toLocaleString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  if (!contact) return null;

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

  return (
    <>
      {store.ui.contactPreviewCardOpen && (
        <div
          data-state={store.ui.contactPreviewCardOpen ? 'open' : 'closed'}
          className='data-[state=open]:animate-slideLeftAndFade data-[state=closed]:animate-slideRightAndFade flex flex-col absolute right-0 -top-[-41px] bottom-0 p-4 max-w-[390px] min-w-[350px] border border-r-0 border-gray-200 z-[1] bg-white'
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
              <IconButton
                size='xs'
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
              onChange={(e) => {
                contact.value.name = e.target.value;
              }}
              onBlur={() => {
                setIsEditName(false);
                contact?.commit();
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
                <span className='mt-2 ml-0.5'> (No org yet)</span>
              )}
            </div>
          )}
          <Input
            size='xs'
            variant='unstyled'
            placeholder='Enter title'
            onFocus={(e) => e.target.select()}
            value={findPrimaryJobRole?.jobTitle || ''}
            className='w-[290px] overflow-hidden text-ellipsis whitespace-nowrap'
            onBlur={() => {
              jobRoleUseCase.submitJobRole(
                String(contactId),
                contact.value.primaryOrganizationId || '',
              );
            }}
            onChange={(e) => {
              const newValue = e.target.value;

              jobRoleUseCase.setJobRole(newValue);

              if (jobRolesStore) {
                jobRolesStore.value.jobTitle = newValue;
              }
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
            <Tags
              placeholder='No tags yet'
              value={
                contact?.value?.tags?.map((tag) => ({
                  value: tag.metadata.id,
                  label: tag.name,
                })) || []
              }
              onChange={(e) => {
                contact.value.tags = e.map(
                  (tag) => store.tags?.value.get(tag.value)?.value,
                ) as Array<Tag>;
                contact.commit();
              }}
            />
          </div>
          <div className='flex flex-col gap-4'>
            <div className='flex gap- items-center w-full '>
              <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500 ' />
                LinkedIn
              </div>
              <div className='flex items-center gap-1 group  '>
                {fromatedUrl ? (
                  <p
                    className='text-sm truncate w-[180px] cursor-default'
                    onClick={() => {
                      copyToClipboard(fromatedUrl, 'LinkedIn profile copied');
                    }}
                  >
                    {fromatedUrl}
                  </p>
                ) : (
                  <p className='text-sm truncate w-[180px] text-gray-400'>
                    LinkedIn profile link
                  </p>
                )}
                {fromatedUrl && (
                  <Link to={href || ''} target='_blank'>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      colorScheme='gray'
                      aria-label='social link'
                      icon={<LinkExternal02 className='text-gray-500' />}
                      className='hover:bg-gray-200 opacity-0 group-hover:opacity-100'
                    />
                  </Link>
                )}
              </div>
            </div>
            <div className='flex gap-1 w-full'>
              <div className='flex items-center gap-2 mr-[41px] text-sm text-gray-500 '>
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
              <div className='flex items-center gap-2 mr-[18px] text-sm text-gray-500'>
                <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                Connected to
              </div>
              <span
                className={cn(
                  'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm',
                  contact?.value?.connectedUsers?.[0]
                    ? 'text-gray-700'
                    : 'text-gray-400',
                )}
              >
                {contact?.value?.connectedUsers?.[0] || 'No one yet'}
              </span>
            </div>
            {contact?.value?.enrichedAt && (
              <div className='text-xs text-gray-500'>
                Last enriched{' '}
                {DateTimeUtils.timeAgo(contact.value.enrichedAt, {
                  addSuffix: true,
                  strict: true,
                  includeMin: true,
                })}
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
