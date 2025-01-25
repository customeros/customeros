import { useRef, useMemo, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

import { useKeyBindings } from 'rooks';
import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';
import { EditJobRole } from '@domain/usecases/contact-preview-card/edit-jobrole.usecase';
import { EditContactNameUseCase } from '@domain/usecases/contact-preview-card/edit-contact-name.usecase';
import { EditContactTagUsecase } from '@domain/usecases/edit-contact-tags-select/edit-contact-tag.usecase.ts';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Tags } from '@shared/components/Tags';
import { TableViewType } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { getTimezone } from '@utils/getTimezone';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { useStore } from '@shared/hooks/useStore';
import { getFormattedLink } from '@utils/getExternalLink';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

import { EmailsSection } from './components';
import { EnrichContactModal } from './components/EnrichContactModal';

export const ContactPreviewCard = observer(() => {
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const linkRef = useRef<HTMLAnchorElement>(null);

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

  const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;
  const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
  const flag = flags[countryA2 || ''];
  const city = contact?.value.locations?.[0]?.locality;
  const timezone = city
    ? cityTimezone.lookupViaCity(city).find((c) => {
        return c.iso2 === contact.value.locations?.[0].countryCodeA2;
      })?.timezone
    : null;

  const linkedInProfile = contact?.value.linkedInAlias;
  const fromatedUrl = getFormattedLink(linkedInProfile || '').replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '',
  );
  const href = contact?.value.linkedInUrl;

  const formatedFollowersCount = contact?.value?.linkedInFollowerCount
    ?.toLocaleString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, ',');

  const jobRoleUseCase = useMemo(
    () => new EditJobRole(String(contactId)),
    [contactId],
  );

  const tagsUseCase = useMemo(
    () => new EditContactTagUsecase(String(contactId)),
    [contactId],
  );
  const contactNameUseCase = useMemo(
    () => new EditContactNameUseCase(String(contactId)),
    [contactId],
  );

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

  const usersStore = store.users
    .toArray()
    .filter((user) => contact.value.connectedUsers?.includes(user.id));

  const users = usersStore.map((user) => user.value.name);

  const lastPositionParams = tabs[contact.value.primaryOrganizationId || ''];
  const hrefOrg = getHref(
    contact.value.primaryOrganizationId || '',
    lastPositionParams || '',
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
              variant='unstyled'
              placeholder='Unknown'
              className='mb-[-8px]'
              value={contact.name || ''}
              onFocus={(e) => e.target.select()}
              onChange={(e) => contact.setName(e.target.value)}
              onBlur={() => {
                contactNameUseCase.execute();
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
                  <Link ref={linkRef} to={hrefOrg || ''}>
                    <p
                      className={cn(
                        'font-medium mt-2 line-clamp-1',
                        company && 'hover:cursor-pointer',
                      )}
                    >
                      {company}
                    </p>
                  </Link>
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
            value={jobRoleUseCase.getJobRole || ''}
            className='w-[290px] overflow-hidden text-ellipsis whitespace-nowrap'
            onChange={(e) => {
              jobRoleUseCase.setJobRole(e.target.value);
            }}
            onBlur={() => {
              jobRoleUseCase.submitJobRole(
                String(contactId),
                contact.value.primaryOrganizationId || '',
              );
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

          <div className='flex justify-between w-full mb-4'>
            <Tags
              dataTest='contact-tags'
              placeholder='No tags yet'
              options={tagsUseCase.tagList}
              onCreate={tagsUseCase.create}
              value={tagsUseCase.selectedTags}
              inputValue={tagsUseCase.searchTerm}
              setInputValue={tagsUseCase.setSearchTerm}
              onChange={(selected) => {
                if (!contact?.value) {
                  throw new Error('Contact store not found');
                }
                tagsUseCase.select(selected?.map((tag) => tag.value));
              }}
              leftAccessory={
                <div className='flex items-center mr-[78px] text-sm text-gray-500'>
                  <Tag01 className='text-gray-500 mr-2' />
                  <span>Tags</span>
                </div>
              }
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
                    {contact?.value.linkedInAlias || fromatedUrl}
                  </p>
                ) : (
                  <p className='text-sm truncate w-[180px] text-gray-400'>
                    LinkedIn profile link
                  </p>
                )}
                {fromatedUrl && (
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    colorScheme='gray'
                    aria-label='social link'
                    icon={<LinkExternal02 className='text-gray-500' />}
                    className='hover:bg-gray-200 opacity-0 group-hover:opacity-100'
                    onClick={() =>
                      window.open(href || '', '_blank', 'noopener')
                    }
                  />
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

              {users.length > 0 ? (
                <span
                  className={cn(
                    'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm text-gray-700',
                  )}
                >
                  {users.join(', ')}
                </span>
              ) : (
                <span
                  className={cn(
                    'overflow-hidden text-ellipsis whitespace-nowrap cursor-not-allowed text-sm text-gray-400',
                  )}
                >
                  {'No one yet'}
                </span>
              )}
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

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=people'}`;
}
