import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';

import { set } from 'lodash';
import { useKeyBindings } from 'rooks';
import cityTimezone from 'city-timezones';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { flags } from '@ui/media/flags';
import { Avatar } from '@ui/media/Avatar';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Mail02 } from '@ui/media/icons/Mail02';
import { getTimezone } from '@utils/getTimezone';
import { useStore } from '@shared/hooks/useStore';
import { Tags } from '@organization/components/Tabs';
import { getFormattedLink } from '@utils/getExternalLink';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import {
  Tag,
  Social,
  TableViewType,
} from '@shared/types/__generated__/graphql.types';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface ContactPreviewCardProps {
  contactId: string;
}

export const ContactPreviewCard = observer(
  ({ contactId }: ContactPreviewCardProps) => {
    const store = useStore();
    const [searchParams] = useSearchParams();
    const [isEditName, setIsEditName] = useState(false);

    const preset = searchParams?.get('preset');
    const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
    const tableType = tableViewDef?.value?.tableType;

    if (tableType !== TableViewType.Contacts && !contactId) {
      store.ui.setContactPreviewCardOpen(false);

      return null;
    }
    const contact = store.contacts.toArray().find((c) => c.id === contactId);

    const fullName = contact?.name || 'Unnamed';
    const src = contact?.value?.profilePhotoUrl;
    const company = contact?.value.organizations?.content?.[0]?.name;
    const role = contact?.value.jobRoles?.[0]?.jobTitle;
    const countryA3 = contact?.value.locations?.[0]?.countryCodeA3;
    const countryA2 = contact?.value.locations?.[0]?.countryCodeA2;
    const flag = flags[countryA2 || ''];
    const city = contact?.value.locations?.[0]?.locality;
    const timezone = city
      ? cityTimezone.lookupViaCity(city)?.[0]?.timezone
      : null;

    const email = contact?.value.emails?.[0]?.email;
    const validationDetails =
      contact?.value.emails?.[0]?.emailValidationDetails;

    const fromatedUrl = contact?.value?.socials?.[0]?.url.replace(
      'https://www.',
      '',
    );

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
        Space: () => {
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
            className='data-[state=open]:animate-slideLeftAndFade data-[state=closed]:animate-slideRightAndFade flex flex-col absolute right-[12px] -top-[-53px] p-4 max-w-[390px] border border-gray-200 rounded-lg z-50 bg-white'
          >
            <Avatar
              size='sm'
              textSize='xs'
              name={fullName}
              variant='circle'
              src={src || undefined}
            />
            <div className='flex items-center gap-1'>
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
                <>
                  <span
                    onClick={() => setIsEditName(true)}
                    className='font-medium mt-2 overflow-hidden text-ellipsis whitespace-nowrap'
                  >
                    {fullName}
                  </span>
                  {company ? (
                    <>
                      <span className='mt-2 text-gray-500'> at </span>
                      <span className='font-medium overflow-hidden text-ellipsis whitespace-nowrap mt-2 '>
                        {company || 'No org yet'}
                      </span>
                    </>
                  ) : (
                    <span className='mt-2'>(No org yet)</span>
                  )}
                </>
              )}
            </div>
            <Input
              size='xs'
              variant='unstyled'
              value={role || ''}
              placeholder='Enter title'
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
                <span className='w-[180px] text-sm'>
                  {getTimezone(timezone || '')} local time
                </span>
              )}
            </div>
            <div className='flex justify-between gap-1 w-full mb-4'>
              <div className='flex items-center gap-2 mr-[70px] text-sm text-gray-500'>
                <Mail02 className='mt-[1px] text-gray-500' />
                Email
              </div>
              <Input
                size='xs'
                variant='unstyled'
                placeholder='Email address'
                className='text-ellipsis ml-[-2px]'
                value={contact?.value.emails?.[0]?.email || ''}
                onBlur={() => {
                  if (!contact?.value?.emails?.[0]?.id) {
                    contact?.addEmail();
                  } else {
                    contact?.updateEmail();
                  }
                }}
                onChange={(e) => {
                  contact?.update(
                    (value) => {
                      set(value, 'emails[0].email', e.target.value);

                      return value;
                    },
                    { mutate: false },
                  );
                }}
              />
              {email && (
                <EmailValidationMessage
                  email={email}
                  validationDetails={validationDetails}
                />
              )}
            </div>

            <div className='flex justify-between gap-1 w-full mb-4'>
              <div className='flex items-center gap-2 mr-[52px] text-sm text-gray-500'>
                <Tag01 className='mt-[1px] text-gray-500' />
                Persona
              </div>
              <Tags
                hideBorder
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

                <Input
                  size='xs'
                  variant='unstyled'
                  value={fromatedUrl}
                  className='text-ellipsis'
                  placeholder='LinkedIn profile link'
                  onChange={(e) => {
                    handleUpdateSocial(e.target.value);
                  }}
                />
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
            </div>
          </div>
        )}
      </>
    );
  },
);