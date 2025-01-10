import { useState } from 'react';
import { Link } from 'react-router-dom';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Spinner } from '@ui/feedback/Spinner';
import { User03 } from '@ui/media/icons/User03';
import { Mail01 } from '@ui/media/icons/Mail01';
import { IconButton } from '@ui/form/IconButton';
import { useEvent } from '@shared/hooks/useEvent';
import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { getFormattedLink } from '@utils/getExternalLink';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
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
  const contactStore = store.contacts.getById(id);

  const handleCreateOption = (value: string) => {
    store.tags?.create(
      { name: value },
      {
        onSucces: (id) => {
          contactStore?.draft();
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
  const formattedLink = getFormattedLink(linkedInProfile || '').replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '/',
  );

  const email =
    contactStore?.value.emails.find((e) => e.primary)?.email ??
    contactStore?.value.emails[0]?.email;

  const jobTitle = contactStore?.value.primaryOrganizationJobRoleTitle;

  const isEnriching = contactStore?.isEnriching;

  if (!contactStore) return null;

  return (
    <>
      <Card
        style={{ paddingBottom: !isExpanded ? '2px' : '10px' }}
        onClick={() => {
          if (!isExpanded) setIsExpanded(true);
        }}
        className={cn(
          isExpanded ? 'bg-white' : 'border-transparent',
          !isExpanded && 'cursor-pointer',
          'px-2 pb-2.5 pt-0.5  max-w-[400px]',
        )}
      >
        <CardHeader style={{ paddingBottom: !isExpanded ? '0' : '8px' }}>
          <div className='flex items-center justify-between w-full'>
            <div className='flex items-center w-full'>
              <div>
                <Avatar
                  size='sm'
                  textSize='sm'
                  variant='outlineCircle'
                  name={contactStore?.value.name ?? ''}
                  className={cn(isEnriching && 'animate-pulse')}
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
                          'cursor-default font-medium text-sm truncate max-w-[200px]',
                          !contactStore.name && 'text-gray-400',
                          !isExpanded && 'cursor-pointer',
                        )}
                      >
                        {isEnriching
                          ? 'Getting name...'
                          : contactStore?.name || 'First & last name'}
                      </span>
                    ) : (
                      <Input
                        size='xxs'
                        variant='unstyled'
                        value={contactStore?.name || ''}
                        dataTest='org-people-contact-name'
                        onFocus={(e) => e.target.select()}
                        onKeyDown={(e) => e.stopPropagation()}
                        className='placeholder:font-medium font-medium min-w-[60px] w-[200px]'
                        onChange={(e) => {
                          contactStore.value.name = e.target.value;
                        }}
                        placeholder={
                          isEnriching ? 'Getting name...' : 'First & last name'
                        }
                        onBlur={() => {
                          contactStore.draft();
                          contactStore.commit();
                        }}
                      />
                    )}

                    <div className='flex items-center h-full gap-1'>
                      {email && !isExpanded && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<Mail01 />}
                          aria-label='send-email'
                          className=' opacity-0 mt-[3px]'
                          onClick={(e) => {
                            dispatchEvent({
                              email: email,
                              openEditor: 'email',
                            });
                            e.stopPropagation();
                          }}
                        />
                      )}
                      {!isExpanded && linkedInProfile && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          icon={<LinkedInSolid02 />}
                          className=' opacity-0 mt-[3px]'
                          aria-label='navigate-to-linkedin'
                          onClick={(e) => {
                            window.open(linkedInProfile, '_blank', 'noopener');
                            e.stopPropagation();
                          }}
                        />
                      )}
                    </div>
                  </div>
                  <div className='flex items-center'>
                    {isEnriching && isExpanded && (
                      <Tooltip
                        open={true}
                        defaultOpen
                        className='z-[9999]'
                        label={`Finding email at ${contactStore.value.primaryOrganizationName}`}
                      >
                        <Spinner
                          size='sm'
                          label='finding email'
                          className='text-gray-400 fill-gray-700 mr-2'
                        />
                      </Tooltip>
                    )}
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      aria-label='collapse'
                      dataTest='org-people-collapse'
                      onClick={() => setIsExpanded(!isExpanded)}
                      className='group-hover/card:opacity-100 opacity-0'
                      icon={
                        !isExpanded ? <ChevronExpand /> : <ChevronCollapse />
                      }
                    />
                    <ContactCardMenu contactId={id} />
                  </div>
                </div>
                {!isExpanded ? (
                  <p
                    className={cn(
                      'text-sm line-clamp-1 cursor-default',
                      !jobTitle && 'text-gray-400',
                      !isExpanded && 'cursor-pointer',
                    )}
                  >
                    {isEnriching
                      ? 'Getting job title...'
                      : jobTitle || 'Job title'}
                  </p>
                ) : (
                  <Input
                    size='xxs'
                    variant='unstyled'
                    onFocus={(e) => e.target.select()}
                    dataTest='org-people-contact-title'
                    onKeyDown={(e) => e.stopPropagation()}
                    placeholder={
                      isEnriching ? 'Getting job title...' : 'Job title'
                    }
                    value={
                      contactStore.value.primaryOrganizationJobRoleTitle || ''
                    }
                    onBlur={() => {
                      contactStore.draft();
                      contactStore.commit();
                    }}
                    onChange={(e) => {
                      set(
                        contactStore.value,
                        'primaryOrganizationJobRoleTitle',
                        e.target.value,
                      );
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

              {linkedInProfile ? (
                <Link
                  target='_blank'
                  to={linkedInProfile || ''}
                  className='cursor-pointer no-underline hover:no-underline'
                >
                  <p
                    className={cn(
                      'text-sm cursor-pointer w-[300px] truncate no-underline hover:no-underline',
                    )}
                  >
                    {formattedLink ?? 'LinkedIn profile URL'}
                  </p>
                </Link>
              ) : (
                <p
                  onClick={() => onOpen()}
                  className={cn(
                    'text-sm cursor-pointer w-[300px] truncate no-underline hover:no-underline',
                    'text-gray-400',
                  )}
                >
                  {'LinkedIn profile URL'}
                </p>
              )}
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
