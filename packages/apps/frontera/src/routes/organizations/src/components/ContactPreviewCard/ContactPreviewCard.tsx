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
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import { Tag, TableViewType } from '@shared/types/__generated__/graphql.types';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

interface ContactPreviewCardProps {
  contactId: string;
}

export const ContactPreviewCard = observer(
  ({ contactId }: ContactPreviewCardProps) => {
    const store = useStore();
    const [searchParams] = useSearchParams();

    const preset = searchParams?.get('preset');
    const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
    const tableType = tableViewDef?.value?.tableType;

    if (tableType !== TableViewType.Contacts) {
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

    useKeyBindings(
      {
        Escape: () => {
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
          <div className='flex flex-col absolute right-[10px] p-4 max-w-[359px] border border-gray-200 rounded-lg z-50 bg-white'>
            <Avatar
              size='xs'
              textSize='xs'
              name={fullName}
              variant='circle'
              src={src || undefined}
            />
            <div className='flex items-center gap-1'>
              <span className='font-semibold'>{fullName}</span> <span>at</span>{' '}
              <span className='font-semibold'>{company}</span>
            </div>
            <p>{role}</p>
            <div className={cn('flex items-center', countryA3 && 'gap-1')}>
              <span>{flag}</span>
              {countryA3 && (
                <span className='ml-2 mt-0.5 overflow-hidden text-ellipsis whitespace-nowrap'>
                  {countryA3}
                </span>
              )}
              {countryA3 && <span>•</span>}
              {city && (
                <span className='overflow-hidden text-ellipsis whitespace-nowrap'>
                  {city}
                </span>
              )}
              {city && <span>•</span>}
              {timezone ? (
                <span className='overflow-hidden text-ellipsis whitespace-nowrap'>
                  {getTimezone(timezone || '')} local time
                </span>
              ) : (
                <span> Invalid timezone</span>
              )}
            </div>
            <div className='flex justify-between gap-1 w-full mb-[-8px]'>
              <div className='flex items-center gap-1 mr-14'>
                <Mail02 className='mt-[1px] text-gray-500' />
                Email
              </div>
              <Input
                variant='unstyled'
                placeholder='Email'
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

            <div className='flex justify-between gap-1 w-full'>
              <div className='flex items-center gap-1 mr-[38px]'>
                <Tag01 className='mt-[1px] text-gray-500' />
                Persona
              </div>
              <Tags
                placeholder='Persona'
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
            <div className='flex flex-col gap-2.5'>
              <div className='flex gap-1 w-full'>
                <div className='flex items-center gap-1 mr-8'>
                  <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                  LinkedIn
                </div>
                <span className='overflow-hidden text-ellipsis whitespace-nowrap '>
                  {fromatedUrl}
                </span>
              </div>
              <div className='flex gap-1 w-full'>
                <div className='flex items-center gap-1 mr-6'>
                  <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                  Followers
                </div>
                <span className='overflow-hidden text-ellipsis whitespace-nowrap '>
                  {formatedFollowersCount}
                </span>
              </div>
              <div className='flex gap-1 w-full'>
                <div className='flex items-center gap-1 mr-8'>
                  <LinkedInSolid02 className='mt-[1px] text-gray-500' />
                  Connected to
                </div>
                <span className='overflow-hidden text-ellipsis whitespace-nowrap '>
                  {contact?.value?.connectedUsers?.[0]?.name}
                </span>
              </div>
            </div>
          </div>
        )}
      </>
    );
  },
);
