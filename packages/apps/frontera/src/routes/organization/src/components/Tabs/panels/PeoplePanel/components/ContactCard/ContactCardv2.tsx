import { useState } from 'react';

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
import { LinkedInSolid } from '@ui/media/icons/LinkedInSolid';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
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
export const ContactCardv2 = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const { dispatchEvent } = useEvent('openEmailEditor');
  const [isCollapsed, setIsCollapsed] = useState(false);
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

  const linkedInProfile = contactStore?.value.socials.find((s) =>
    s.url.includes('linkedin'),
  )?.url;

  const email =
    contactStore?.value.emails.find((e) => e.primary)?.email ??
    contactStore?.value.emails[0]?.email;

  if (!contactStore) return null;

  return (
    <>
      <Card
        style={{ paddingBottom: !isCollapsed ? '0' : '10px' }}
        className={cn(
          isCollapsed ? 'bg-white' : 'border-transparent',
          'px-2 pb-2.5 pt-0.5 group-hover/card:border-gray-200 group-hover/card:bg-white',
        )}
      >
        <CardHeader style={{ paddingBottom: !isCollapsed ? '0' : '8px' }}>
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
                <div className='flex justify-between group/action-buttons'>
                  <div className='flex justify-start gap-4'>
                    {!isCollapsed ? (
                      <p className='cursor-default font-medium'>
                        {contactStore?.name ?? ''}
                      </p>
                    ) : (
                      <Input
                        size='xxs'
                        variant='unstyled'
                        placeholder='First & last name'
                        value={contactStore?.name ?? ''}
                        dataTest='org-people-contact-name'
                        onBlur={() => contactStore.commit()}
                        className='placeholder:font-medium font-medium min-w-[60px] max-w-fit'
                        onChange={(e) => {
                          contactStore.value.name = e.target.value;
                        }}
                      />
                    )}

                    <div className='flex gap-1'>
                      {email && !isCollapsed && (
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
                      {!isCollapsed && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<LinkedInSolid />}
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
                      onClick={() => setIsCollapsed(!isCollapsed)}
                    />
                    <ContactCardMenu contactId={id} />
                  </div>
                </div>
                {!isCollapsed ? (
                  <p className='text-sm line-clamp-1 cursor-default'>
                    {contactStore?.value?.jobRoles?.[0]?.jobTitle ?? ''}
                  </p>
                ) : (
                  <Input
                    size='xxs'
                    variant='unstyled'
                    placeholder='Job title'
                    value={contactStore?.value?.jobRoles?.[0]?.jobTitle ?? ''}
                    onBlur={() => {
                      contactStore.commit();
                    }}
                    onChange={(e) => {
                      contactStore.draft();

                      contactStore.value.jobRoles[0].jobTitle = e.target.value;
                    }}
                  />
                )}
              </div>
            </div>
          </div>
        </CardHeader>
        {isCollapsed && (
          <CardContent className='pl-2 pr-0 pb-0 gap-1 flex flex-col'>
            <ContactLocation contactId={id} />
            <ContactJobExperience contactId={id} />
            <EmailsSection contactId={id} />

            <div>
              <Linkedin className='text-gray-500 mr-4' />
              <span
                onClick={() => onOpen()}
                className={cn('text-sm', !linkedInProfile && 'text-gray-400')}
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
            <CardFooter className='pt-0.5 pb-0.5 px-0'>
              <span className='text-[12px] text-grayModern-500'>
                {contactStore?.value?.updatedAt && updatedDaysAgo}
              </span>
            </CardFooter>
          </CardContent>
        )}
      </Card>
      <AddLinkedInToContactModal open={open} onClose={onClose} />
    </>
  );
});
