import { useState } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Tag01 } from '@ui/media/icons/Tag01';
import { User03 } from '@ui/media/icons/User03';
import { Mail01 } from '@ui/media/icons/Mail01';
import { IconButton } from '@ui/form/IconButton';
import { useEvent } from '@shared/hooks/useEvent';
import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import {
  Tag,
  DataSource,
  EntityType,
} from '@shared/types/__generated__/graphql.types';
import {
  Card,
  CardHeader,
  CardFooter,
  CardContent,
} from '@ui/presentation/Card/Card';

import { EmailsSection } from './components/EmailsSection';
import { Tags } from '../../../AboutPanel/components/tags';
import { ContactCardMenu, ContactLocation } from './components';
import { ContactJobExperience } from './components/ContactJobExperience';
import { AddLinkedInToContactModal } from './components/AddLinkedInToContactModal';

interface ContactCardProps {
  id: string;
}
export const ContactCard = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const { dispatchEvent } = useEvent('openEmailEditor');
  const [isExpanded, setIsExpanded] = useState(false);
  const { onOpen, onClose, open } = useDisclosure();
  const contactStore = store.contacts.value.get(id);

  const handleCreateOption = (value: string) => {
    store.tags?.create(
      { name: value },
      {
        onSucces: (id) => {
          contactStore?.value.tags?.push({
            name: value,
            metadata: {
              id,
              source: DataSource.Openline,
              sourceOfTruth: DataSource.Openline,
              appSource: 'organization',
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
            },
            entityType: EntityType.Contact,
          } as Tag);
          contactStore?.commit();
        },
      },
    );
  };

  const updatedDaysAgo =
    DateTimeUtils.getDaysSinceDate(contactStore?.value.updatedAt) === 1
      ? `Last changed ${DateTimeUtils.getDaysSinceDate(
          contactStore?.value.updatedAt,
        )} day ago`
      : `Last changed ${DateTimeUtils.getDaysSinceDate(
          contactStore?.value.updatedAt,
        )} days ago`;

  const linkedInProfile = contactStore?.value.linkedInUrl;

  const email =
    contactStore?.value.emails.find((e) => e.primary)?.email ??
    contactStore?.value.emails[0]?.email;

  const jobTitle = contactStore?.value.primaryOrganizationJobRoleTitle;

  if (!contactStore) return null;

  return (
    <>
      <Card
        style={{ paddingBottom: !isExpanded ? '0' : '10px' }}
        className={cn(
          isExpanded ? 'bg-white' : 'border-transparent',
          'px-2 pb-2.5 pt-0.5 group-hover/card:border-gray-200 group-hover/card:bg-white',
        )}
      >
        <CardHeader style={{ paddingBottom: !isExpanded ? '0' : '8px' }}>
          <div className='flex items-center justify-between w-full'>
            <div className='flex items-center w-full'>
              <div>
                <Avatar
                  size='sm'
                  variant='outlineCircle'
                  name={contactStore?.value.name ?? ''}
                  icon={<User03 className='text-gray-700 size-6' />}
                  src={
                    contactStore?.value?.profilePhotoUrl
                      ? contactStore.value.profilePhotoUrl
                      : undefined
                  }
                />
              </div>

              <div className='flex flex-col w-full ml-2'>
                <div className='flex justify-between group/action-buttons min-w-max'>
                  <div className='flex justify-start gap-4 h-full '>
                    {!isExpanded ? (
                      <span
                        className={cn(
                          'cursor-default font-medium text-sm',
                          !contactStore.name && 'text-gray-400',
                        )}
                      >
                        {contactStore?.name || 'First & last name'}
                      </span>
                    ) : (
                      <Input
                        size='xxs'
                        variant='unstyled'
                        placeholder='First & last name'
                        value={contactStore?.name ?? ''}
                        dataTest='org-people-contact-name'
                        className='placeholder:font-medium font-medium min-w-[60px]'
                        onChange={(e) => {
                          contactStore.value.name = e.target.value;
                        }}
                        onBlur={() => {
                          contactStore.draft();
                          contactStore.commit();
                        }}
                      />
                    )}

                    <div className='flex h-full gap-1'>
                      {email && !isExpanded && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<Mail01 />}
                          aria-label='send-email'
                          className='group-hover/action-buttons:opacity-100 opacity-0'
                          onClick={() =>
                            dispatchEvent({ email: email, openEditor: 'email' })
                          }
                        />
                      )}
                      {!isExpanded && linkedInProfile && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<LinkedInSolid02 />}
                          aria-label='navigate-to-linkedin'
                          className='group-hover/action-buttons:opacity-100 opacity-0'
                          onClick={() =>
                            window.open(linkedInProfile, '_blank', 'noopener')
                          }
                        />
                      )}
                    </div>
                  </div>
                  <div>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      aria-label='collapse'
                      icon={<ChevronCollapse />}
                      dataTest='org-people-collapse'
                      onClick={() => setIsExpanded(!isExpanded)}
                      className='group-hover/card:opacity-100 opacity-0'
                    />
                    <ContactCardMenu contactId={id} />
                  </div>
                </div>
                {!isExpanded ? (
                  <p
                    className={cn(
                      'text-sm line-clamp-1 cursor-default',
                      !jobTitle && 'text-gray-400',
                    )}
                  >
                    {jobTitle || 'Job title'}
                  </p>
                ) : (
                  <Input
                    size='xxs'
                    variant='unstyled'
                    placeholder='Job title'
                    dataTest='org-people-contact-title'
                    value={
                      contactStore.value.primaryOrganizationJobRoleTitle || ''
                    }
                    onBlur={() => {
                      contactStore.draft();
                      contactStore.commit();
                    }}
                    onChange={(e) => {
                      if (e.target.value !== '') {
                        set(
                          contactStore.value,
                          'primaryOrganizationJobRoleTitle',
                          e.target.value,
                        );
                      }
                    }}
                  />
                )}
              </div>
            </div>
          </div>
        </CardHeader>
        {isExpanded && (
          <CardContent className='pl-2 pr-0 pb-0 gap-3 flex flex-col justify-center'>
            {contactStore.value.locations?.length > 0 && (
              <ContactLocation contactId={id} />
            )}
            {contactStore.value.primaryOrganizationJobRoleStartDate && (
              <ContactJobExperience contactId={id} />
            )}
            <EmailsSection contactId={id} />

            <div className='flex items-center max-h-6'>
              <Linkedin className='text-gray-500 mr-4' />
              <span
                data-test={'org-people-linkedin'}
                className={cn(
                  'text-sm cursor-pointer',
                  !linkedInProfile && 'text-gray-400',
                )}
                onClick={() => {
                  if (!linkedInProfile) {
                    onOpen();
                  } else {
                    window.open(linkedInProfile, '_blank', 'noopener');
                  }
                }}
              >
                {linkedInProfile || 'LinkedIn profile URL'}
              </span>
            </div>

            <Tags
              placeholder='Tags'
              dataTest='org-about-tags'
              className='min-h-4 text-sm'
              inputPlaceholder='Search...'
              onCreate={handleCreateOption}
              leftAccessory={<Tag01 className='mr-4 text-gray-500 size-4' />}
              value={
                contactStore.value.tags?.map((t) => ({
                  value: t.metadata.id,
                  label: t.name,
                })) ?? []
              }
              options={store.tags
                .getByEntityType(EntityType.Contact)
                .map((t) => ({
                  value: t.id,
                  label: t.value?.name,
                }))}
              onChange={(selection) => {
                const tags = selection
                  .map((o) => store.tags.getById(o.value)?.value)
                  .filter(Boolean);

                contactStore.value.tags = tags as TagDatum[];
                contactStore.draft();
                contactStore.commit();
              }}
            />
            <CardFooter className='pt-0.5 pb-0.5 px-0 max-h-5'>
              <span className='text-[12px] text-grayModern-500'>
                {contactStore?.value?.updatedAt && updatedDaysAgo}
              </span>
            </CardFooter>
          </CardContent>
        )}
      </Card>
      <AddLinkedInToContactModal open={open} contactId={id} onClose={onClose} />
    </>
  );
});
