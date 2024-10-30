import React, { useRef, useMemo, useEffect, MouseEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';
import { formatDistanceToNow } from 'date-fns/formatDistanceToNow';
import { differenceInCalendarMonths } from 'date-fns/differenceInCalendarMonths';

import { cn } from '@ui/utils/cn';
import { Select } from '@ui/form/Select';
import { Input } from '@ui/form/Input/Input';
import { Clock } from '@ui/media/icons/Clock';
import { Check } from '@ui/media/icons/Check';
import { Mail01 } from '@ui/media/icons/Mail01';
import { User03 } from '@ui/media/icons/User03';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { Share07 } from '@ui/media/icons/Share07';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useStore } from '@shared/hooks/useStore';
import { Users01 } from '@ui/media/icons/Users01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Tags } from '@organization/components/Tabs/shared/';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Tag, Social, Contact, DataSource } from '@graphql/types';
import { PhoneOutgoing02 } from '@ui/media/icons/PhoneOutgoing02';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { useContactCardMeta } from '@organization/state/ContactCardMeta.atom';
import { SocialIconInput } from '@organization/components/Tabs/shared/SocialIconInput';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '@ui/form/InputGroup/InputGroup';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';
import { EmailValidationMessage } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

import { timezoneOptions } from '../util';
import { TimezoneSelect } from './TimezoneSelect';

const roleOptions = [
  {
    value: 'Decision Maker',
    label: 'Decision Maker',
  },
  {
    value: 'Influencer',
    label: 'Influencer',
  },
  {
    value: 'User',
    label: 'User',
  },
  {
    value: 'Stakeholder',
    label: 'Stakeholder',
  },
  {
    value: 'Gatekeeper',
    label: 'Gatekeeper',
  },
  {
    value: 'Champion',
    label: 'Champion',
  },
  {
    value: 'Data Owner',
    label: 'Data Owner',
  },
];

interface ContactCardProps {
  id: string;
  contact: Contact;
  organizationName?: string;
}

export const ContactCard = observer(
  ({ id, organizationName }: ContactCardProps) => {
    const store = useStore();
    const cardRef = useRef<HTMLDivElement>(null);
    const [{ expandedId, initialFocusedField }, setExpandedCardId] =
      useContactCardMeta();
    const isExpanded = expandedId === id;
    const { open: isOpen, onOpen, onClose } = useDisclosure();

    useOutsideClick({
      ref: cardRef,
      handler: () => {
        if (expandedId === id) {
          setExpandedCardId({
            expandedId: undefined,
            initialFocusedField: null,
          });
        }
      },
    });

    const contactStore = store.contacts.value.get(id);
    const emailInputRef = useRef<HTMLInputElement | null>(null);
    const nameInputRef = useRef<HTMLInputElement | null>(null);
    const oldEmail = useMemo(
      () => contactStore?.value.emails?.[0]?.email,
      [contactStore?.isLoading],
    );

    const toggle = (e: MouseEvent<HTMLDivElement>) => {
      if (
        ['name', 'role', 'title'].includes((e.target as HTMLDivElement)?.id)
      ) {
        setExpandedCardId({
          expandedId: id,
          initialFocusedField: null,
        });

        return;
      }

      if (isExpanded) {
        setExpandedCardId({ expandedId: undefined, initialFocusedField: null });
      } else {
        setExpandedCardId({
          expandedId: id,
          initialFocusedField: null,
        });
      }
    };

    useEffect(() => {
      if (expandedId === id && initialFocusedField) {
        if (initialFocusedField === 'name') {
          nameInputRef.current?.focus();

          return;
        }

        if (initialFocusedField === 'email') {
          emailInputRef.current?.focus();

          return;
        }
      }
    }, [expandedId, initialFocusedField, emailInputRef]);

    const timeAt = (() => {
      const startedAt = contactStore?.value?.jobRoles?.[0]?.startedAt;

      if (!startedAt) return undefined;

      const months = Math.abs(
        differenceInCalendarMonths(new Date(startedAt), new Date()),
      );

      if (months < 0) return `Less than a month at ${organizationName}`;
      if (months === 1) return `${months} month at ${organizationName}`;
      if (months > 1 && months < 12)
        return `${months} months at ${organizationName}`;
      if (months === 12) return `1 year at ${organizationName}`;
      if (months > 12)
        return `${formatDistanceToNow(
          new Date(startedAt),
        )} at ${organizationName}`;
    })();

    const handleDelete = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      store.contacts.softDelete(id);
      onClose();
    };

    const toggleConfirmDelete = (e: MouseEvent) => {
      e.stopPropagation();
      e.preventDefault();
      onOpen();
    };

    const handleFindEmail = () => {
      contactStore?.findEmail();
    };

    const handleChange = (
      e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
      contactStore?.update((value) => {
        const property = e.target.name as keyof Contact;

        if (!property) return value;

        value[property] = e.target.value;

        return value;
      });
    };

    const handleCreateOption = (value: string) => {
      store.tags?.create({ name: value });

      contactStore?.update((org) => {
        org.tags = [
          ...(org.tags || []),
          {
            id: value,
            name: value,
            metadata: {
              id: value,
              source: DataSource.Openline,
              sourceOfTruth: DataSource.Openline,
              appSource: 'organization',
              created: new Date().toISOString(),
              lastUpdated: new Date().toISOString(),
            },
            appSource: 'organization',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            source: DataSource.Openline,
          },
        ];

        return org;
      });
    };

    const enrichedContact = contactStore?.value.enrichDetails;
    const enrichingStatus =
      !enrichedContact?.enrichedAt &&
      enrichedContact?.requestedAt &&
      !enrichedContact?.failedAt;

    const enrichedEmailStatus =
      !enrichedContact?.emailEnrichedAt &&
      enrichedContact?.emailRequestedAt &&
      !enrichedContact?.emailFound;

    return (
      <>
        <Card
          ref={cardRef}
          key={contactStore?.getId()}
          className={cn(
            'bg-white w-full group rounded-lg border-[1px] border-gray-200 cursor-pointer hover:shadow-md ',
            isExpanded ? 'shadow-md' : 'shadow-xs',
            'ease-linear',
            'transition-all',
            'duration-1000',
          )}
        >
          <CardHeader onClick={toggle} className={cn('flex p-4 relative')}>
            <div className='flex flex-col'>
              {enrichingStatus && (
                <div className='flex items-center justify-start gap-2 border-[1px] text-sm border-grayModern-100 bg-grayModern-50 rounded-[4px] py-1 px-2 mb-4 w-full'>
                  <Spinner
                    label='enriching org'
                    className='text-grayModern-300 fill-grayModern-500 size-4'
                  />
                  <span className='font-medium'>
                    We're enriching this contact's details...
                  </span>
                </div>
              )}
              <div className='flex'>
                <Avatar
                  variant='shadowed'
                  name={contactStore?.value.name ?? ''}
                  icon={<User03 className='text-primary-700 size-6' />}
                  src={
                    contactStore?.value?.profilePhotoUrl
                      ? contactStore.value.profilePhotoUrl
                      : undefined
                  }
                />

                <div className='ml-4 flex flex-col flex-1'>
                  <Input
                    size='xs'
                    name='name'
                    ref={nameInputRef}
                    placeholder='Name'
                    onChange={handleChange}
                    value={contactStore?.name ?? ''}
                    dataTest='org-people-contact-name'
                    className='font-semibold text-gray-700'
                  />
                  <Input
                    size='xs'
                    name='prefix'
                    placeholder='Title'
                    className='text-gray-500'
                    dataTest='org-people-contact-title'
                    value={contactStore?.value?.jobRoles?.[0]?.jobTitle ?? ''}
                    onChange={(e) => {
                      contactStore?.update((value) => {
                        set(value, 'jobRoles[0].jobTitle', e.target.value);

                        return value;
                      });
                    }}
                  />
                  <Select
                    isMulti
                    size='xs'
                    name='role'
                    options={roleOptions}
                    placeholder='Choose job roles'
                    dataTest='org-people-contact-job-roles'
                    value={
                      contactStore?.value?.jobRoles?.[0]?.description
                        ?.split(',')
                        .filter(Boolean)
                        .map((v) => ({ value: v, label: v })) ?? []
                    }
                    onChange={(opt) => {
                      contactStore?.update((value) => {
                        const description = opt
                          .map((v: SelectOption) => v.value)
                          .join(',');

                        set(value, 'jobRoles[0].description', description);

                        return value;
                      });
                    }}
                  />
                </div>
              </div>
            </div>
            {isExpanded && (
              <IconButton
                size='xs'
                variant='ghost'
                onClick={onClose}
                colorScheme='gray'
                aria-label='Close'
                id='collapse-button'
                dataTest='org-people-contact-close'
                icon={<Check className='text-gray-500' />}
                className='absolute z-50 top-2 right-2 p-1 opacity-0 pointer-events-auto transition-opacity duration-300 group-hover:opacity-100 "'
              />
            )}

            {!isExpanded && (
              <IconButton
                size='sm'
                variant='ghost'
                colorScheme='gray'
                id='confirm-button'
                aria-label='Delete contact'
                onClick={toggleConfirmDelete}
                dataTest='org-people-contact-delete'
                icon={<Trash01 className='text-gray-400' />}
                className='hover:bg-error-100 *:hover:text-error-500 absolute z-50 top-2 right-2 p-1 opacity-0 pointer-events-auto transition-opacity duration-300 group-hover:opacity-100 "'
              />
            )}
          </CardHeader>
          {isExpanded && (
            <CardContent
              className={cn('flex flex-col', isExpanded ? 'h-auto' : 'h-0')}
            >
              <InputGroup>
                <LeftElement>
                  <Tooltip hasArrow label='Click to autopopulate'>
                    <span>
                      {enrichedEmailStatus ? (
                        <Spinner
                          size='sm'
                          label='Finding email'
                          className='text-gray-300 fill-gray-500'
                        />
                      ) : (
                        <Mail01
                          onClick={handleFindEmail}
                          className='text-gray-500 hover:text-gray-700 transition-colors'
                        />
                      )}
                    </span>
                  </Tooltip>
                </LeftElement>
                <Input
                  variant='unstyled'
                  ref={emailInputRef}
                  placeholder='Email'
                  dataTest='org-people-contact-email'
                  value={contactStore?.value?.emails?.[0]?.email ?? ''}
                  onBlur={() => {
                    contactStore?.updateEmail(oldEmail ?? '');
                  }}
                  onChange={(e) => {
                    contactStore?.update(
                      (value) => {
                        set(value, 'emails[0].email', e.target.value);

                        return value;
                      },
                      { mutate: false },
                    );
                  }}
                />
                <RightElement>
                  <EmailValidationMessage
                    email={contactStore?.value?.emails?.[0]?.email ?? ''}
                    validationDetails={
                      contactStore?.value?.emails?.[0]?.emailValidationDetails
                    }
                  />
                </RightElement>
              </InputGroup>

              <InputGroup>
                <LeftElement>
                  <PhoneOutgoing02 className='text-gray-500' />
                </LeftElement>
                <Input
                  variant='unstyled'
                  placeholder='Phone number'
                  dataTest='org-people-contact-phone-number'
                  value={
                    contactStore?.value.phoneNumbers?.[0]?.rawPhoneNumber ?? ''
                  }
                  onBlur={() => {
                    if (!contactStore?.value.phoneNumbers?.[0]?.id) {
                      contactStore?.addPhoneNumber();
                    } else {
                      contactStore?.updatePhoneNumber();
                    }
                  }}
                  onChange={(e) => {
                    contactStore?.update(
                      (value) => {
                        set(
                          value,
                          'phoneNumbers[0].rawPhoneNumber',
                          e.target.value,
                        );

                        return value;
                      },
                      { mutate: false },
                    );
                  }}
                />
              </InputGroup>

              {/* TODO: replace with FormInput. currently displayed as a text just for demoing purposes */}
              {timeAt && (
                <div className='flex items-center h-[39px]'>
                  <Calendar className='text-gray-500' />
                  <p className='ml-[14px] cursor-text capitalize'>{timeAt}</p>
                </div>
              )}
              {/* END TODO */}

              <Tags
                placeholder='Personas'
                onCreateOption={handleCreateOption}
                dataTest='org-people-contact-personas'
                icon={
                  <Users01 className='text-gray-500 w-[18px] h-4 mr-[10px] mt-[6px] ' />
                }
                value={
                  contactStore?.value?.tags?.map((t) => ({
                    label: t.name,
                    value: t.id,
                  })) ?? []
                }
                onChange={(e) => {
                  contactStore?.update((c) => {
                    c.tags =
                      (e
                        .map((tag) => store.tags?.value.get(tag.value)?.value)
                        .filter(Boolean) as Array<Tag>) ?? [];

                    return c;
                  });
                }}
              />

              <SocialIconInput
                placeholder='Social link'
                dataTest='org-people-contact-social-link'
                leftElement={<Share07 className='text-gray-500' />}
                value={
                  contactStore?.value?.socials?.map((s) => ({
                    label: s?.alias ? `linkedin.com/in/${s.alias}` : s.url,
                    value: s.id,
                  })) ?? []
                }
                onCreate={(value) => {
                  contactStore?.update((prev) => {
                    prev.socials.push({
                      id: crypto.randomUUID(),
                      url: value,
                    } as Social);

                    return prev;
                  });
                }}
                onChange={(e) => {
                  const id = e.target.id;

                  contactStore?.update((value) => {
                    const foundIndex = value.socials.findIndex(
                      (s) => s.id === id,
                    );

                    if (foundIndex !== -1) {
                      value.socials[foundIndex].url = e.target.value;
                      set(
                        value,
                        ['socials', foundIndex, 'url'],
                        e.target.value,
                      );
                    }

                    return value;
                  });
                }}
              />
              <TimezoneSelect
                isClearable
                placeholder='Timezone'
                options={timezoneOptions}
                dataTest='org-people-contact-timezone'
                leftElement={<Clock className='text-gray-500 mr-3' />}
                value={timezoneOptions.find(
                  (v) => v.value === contactStore?.value?.timezone,
                )}
                onChange={(opt) => {
                  contactStore?.update((value) => {
                    value.timezone = opt?.value;

                    return value;
                  });
                }}
              />
              {/* <AutoresizeTextarea
                className='items-start'
                name='description'
                placeholder='Notes'
                onChange={handleChange}
                value={contactStore?.value?.description ?? ''}
                leftElement={<File02 className='text-gray-500 mt-1 mr-1' />}
              /> */}
            </CardContent>
          )}
        </Card>
        <ConfirmDeleteDialog
          isOpen={isOpen}
          hideCloseButton
          onClose={onClose}
          onConfirm={handleDelete}
          label='Delete this contact?'
          confirmButtonLabel='Delete contact'
        />
      </>
    );
  },
);
