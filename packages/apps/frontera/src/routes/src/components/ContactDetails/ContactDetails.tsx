import { useParams } from 'react-router-dom';
import { useMemo, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { ContactDetailsUseCase } from '@domain/usecases/contact-details/contact-details.usecase';
import { EditContactTagUsecase } from '@domain/usecases/edit-contact-tags-select/edit-contact-tag.usecase';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Avatar } from '@ui/media/Avatar';
import { DateTimeUtils } from '@utils/date';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Spinner } from '@ui/feedback/Spinner';
import { User03 } from '@ui/media/icons/User03';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Copy01 } from '@ui/media/icons/Copy01';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { useEvent } from '@shared/hooks/useEvent';
import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Building07 } from '@ui/media/icons/Building07';
import { getFormattedLink } from '@utils/getExternalLink';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { LinkedInSolid02 } from '@ui/media/icons/LinkedInSolid02';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import {
  Tag,
  DataSource,
  EntityType,
} from '@shared/types/__generated__/graphql.types';

import { Tags } from '../Tags';
import { EmailsSection } from './components/EmailsSection';
import { ContactLocation } from './components/ContactLocation';
import { ContactCardMenu } from './components/ContactCardMenu';
import { ContactJobExperience } from './components/ContactJobExperience';
import { AddLinkedInToContactModal } from './components/AddLinkedInToContactModal';

interface ContactDetailsProps {
  id: string;
  expandAll?: boolean;
  isExpandble?: boolean;
}

export const ContactDetails = observer(
  ({ id, isExpandble = false, expandAll = false }: ContactDetailsProps) => {
    const store = useStore();

    const [isExpanded, setIsExpanded] = useState(false);
    const { dispatchEvent } = useEvent('openEmailEditor');
    const { onOpen, onClose, open } = useDisclosure();
    const orgId = useParams()?.id as string;
    const contactStore = store.contacts.getById(id);
    const jobRoles = contactStore?.jobRoles;

    const [_, copyToClipboard] = useCopyToClipboard();

    const findPrimaryJobRole = jobRoles?.find(
      (j) => j.primary && j.contact?.metadata.id === id,
    );
    const tagsUseCase = useMemo(
      () => new EditContactTagUsecase(String(id)),
      [id],
    );

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

    useEffect(() => {
      if (expandAll || !isExpandble) {
        setIsExpanded(true);
      } else {
        setIsExpanded(false);
      }
    }, [expandAll, isExpandble]);

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

    const isEnriching = contactStore?.isEnriching;

    const contactDetails = useMemo(() => new ContactDetailsUseCase(id), [id]);

    if (!contactStore) return null;

    return (
      <>
        <div
          style={{ paddingBottom: !isExpanded && isExpandble ? '2px' : '10px' }}
          onClick={() => {
            if (!isExpandble) return;
            if (!isExpanded && isExpandble) setIsExpanded(true);
          }}
          className={cn(
            isExpanded && isExpandble ? 'bg-white' : 'border-transparent',
            !isExpanded && isExpandble && 'cursor-pointer',
            isExpandble
              ? 'px-2 pb-2.5 pt-0.5 group-hover/card:border-gray-200 group-hover/card:bg-white max-w-[400px] rounded-lg border bg-card text-card-foreground shadow'
              : 'p-4',
          )}
        >
          <div
            style={{ paddingBottom: !isExpanded && isExpandble ? '0' : '8px' }}
          >
            {!isExpandble && (
              <div className='flex items-center justify-between pb-3'>
                <p className='text-base font-medium '>Contact overview</p>
                <IconButton
                  size='xs'
                  variant='ghost'
                  icon={<XClose />}
                  aria-label='close-contact-overview'
                  onClick={() => store.ui.setShowPreviewCard(false)}
                />
              </div>
            )}
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
                      {!isExpanded && isExpandble ? (
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
                          value={contactStore.name || ''}
                          dataTest='org-people-contact-name'
                          onFocus={(e) => e.target.select()}
                          onKeyDown={(e) => e.stopPropagation()}
                          className='placeholder:font-medium font-medium min-w-[60px] w-[200px]'
                          onBlur={() => {
                            contactDetails.changeContactName();
                          }}
                          onChange={(e) =>
                            contactStore.setName(e.target.value || '')
                          }
                          placeholder={
                            isEnriching
                              ? 'Getting name...'
                              : 'First & last name'
                          }
                        />
                      )}

                      <div className='flex items-center h-full gap-1'>
                        {email && !isExpanded && isExpandble && (
                          <IconButton
                            size='xxs'
                            variant='ghost'
                            icon={<Mail01 />}
                            aria-label='send-email'
                            className='opacity-0 mt-[3px]'
                            onClick={(e) => {
                              dispatchEvent({
                                email: email,
                                openEditor: 'email',
                              });
                              e.stopPropagation();
                            }}
                          />
                        )}
                        {!isExpanded && linkedInProfile && isExpandble && (
                          <IconButton
                            size='xxs'
                            variant='ghost'
                            icon={<LinkedInSolid02 />}
                            className='opacity-0 mt-[3px]'
                            aria-label='navigate-to-linkedin'
                            onClick={(e) => {
                              window.open(
                                linkedInProfile,
                                '_blank',
                                'noopener',
                              );
                              e.stopPropagation();
                            }}
                          />
                        )}
                      </div>
                    </div>
                    <div className='flex items-center'>
                      {isEnriching && isExpanded && isExpandble && (
                        <Tooltip
                          className='z-[9999]'
                          label={`Finding email at ${contactStore.value.primaryOrganizationName}`}
                        >
                          <Spinner
                            size='sm'
                            label='finding email'
                            className='text-gray-400 fill-gray-700 mr-2 group-hover/card:opacity-100 opacity-0'
                          />
                        </Tooltip>
                      )}
                      {isExpandble && (
                        <IconButton
                          size='xxs'
                          variant='ghost'
                          aria-label='collapse'
                          dataTest='org-people-collapse'
                          onClick={() => setIsExpanded(!isExpanded)}
                          className='group-hover/card:opacity-100 opacity-0'
                          icon={
                            !isExpanded ? (
                              <ChevronExpand />
                            ) : (
                              <ChevronCollapse />
                            )
                          }
                        />
                      )}
                      <ContactCardMenu contactId={id} />
                    </div>
                  </div>
                  {!isExpanded && isExpandble ? (
                    <p
                      className={cn(
                        'text-sm line-clamp-1 cursor-default',
                        !findPrimaryJobRole?.jobTitle && 'text-gray-400',
                        !isExpanded && 'cursor-pointer',
                      )}
                    >
                      {isEnriching
                        ? 'Getting job title...'
                        : findPrimaryJobRole?.jobTitle || 'No job title yet'}
                    </p>
                  ) : (
                    <Input
                      size='xxs'
                      variant='unstyled'
                      onFocus={(e) => e.target.select()}
                      dataTest='org-people-contact-title'
                      onKeyDown={(e) => e.stopPropagation()}
                      value={contactDetails.getJobRole ?? ''}
                      onBlur={() => {
                        contactDetails.submitJobRole(id, orgId);
                      }}
                      onChange={(e) => {
                        contactDetails.setJobRole(e.target.value);
                      }}
                      placeholder={
                        isEnriching
                          ? 'Getting job title...'
                          : 'No job title yet'
                      }
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
          {(isExpanded || !isExpandble) && (
            <div className='pl-2 pr-0 pb-0 gap-3 flex flex-col justify-center'>
              {contactStore.value.locations?.length > 0 && (
                <ContactLocation contactId={id} />
              )}
              {!isExpandble && contactStore.value.primaryOrganizationName && (
                <div className='flex items-center'>
                  <Building07 className='text-gray-500' />
                  <span className='text-sm  ml-4'>
                    {contactStore.value.primaryOrganizationName}
                  </span>
                </div>
              )}
              {contactStore.value.primaryOrganizationJobRoleStartDate && (
                <ContactJobExperience contactId={id} />
              )}
              <EmailsSection contactId={id} />

              <div className='flex items-center max-h-6 group/linkedin'>
                <Linkedin className='text-gray-500 mr-4' />
                <div className='flex items-center w-full '>
                  {linkedInProfile ? (
                    <span
                      onClick={() =>
                        window.open(linkedInProfile, '_blank', 'noopener')
                      }
                      className={cn(
                        'text-sm cursor-pointer max-w-[300px] truncate no-underline hover:underline',
                      )}
                    >
                      {contactStore.value.linkedInAlias &&
                      contactStore.value.linkedInAlias?.length > 1
                        ? contactStore.value.linkedInAlias
                        : formattedLink}
                    </span>
                  ) : (
                    <span
                      onClick={() => onOpen()}
                      data-test='org-people-linkedin'
                      className={cn(
                        'text-sm cursor-pointer max-w-[300px] truncate no-underline hover:no-underline',
                        'text-gray-400',
                      )}
                    >
                      {'LinkedIn profile URL'}
                    </span>
                  )}
                  {linkedInProfile && (
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      title='Copy link'
                      colorScheme='gray'
                      aria-label='social link'
                      icon={<Copy01 className='text-gray-500' />}
                      className='hover:bg-gray-200 ml-2 opacity-0 group-hover/linkedin:opacity-100'
                      onClick={() =>
                        copyToClipboard(
                          contactStore.value.linkedInUrl || '',
                          'LinkedIn profile copied',
                        )
                      }
                    />
                  )}
                </div>
              </div>
              <Tags
                placeholder='Tags'
                dataTest='contact-tags'
                className='min-h-4 text-sm'
                inputPlaceholder='Search...'
                onCreate={handleCreateOption}
                options={tagsUseCase.tagList}
                value={tagsUseCase.selectedTags}
                inputValue={tagsUseCase.searchTerm}
                setInputValue={tagsUseCase.setSearchTerm}
                leftAccessory={<Tag01 className='mr-4 text-gray-500 size-4' />}
                onChange={(selected) => {
                  if (Array.isArray(selected)) {
                    if (!contactStore?.value) {
                      throw new Error('Contact store not found');
                    }
                    tagsUseCase.select(selected?.map((tag) => tag.value));
                  }
                }}
              />

              <div className='pt-0.5 pb-0.5 px-0 max-h-5'>
                <span className='text-[12px] text-grayModern-500'>
                  {contactStore?.value?.updatedAt && updatedDaysAgo}
                </span>
              </div>
            </div>
          )}
        </div>
        <AddLinkedInToContactModal
          open={open}
          contactId={id}
          onClose={onClose}
        />
      </>
    );
  },
);
