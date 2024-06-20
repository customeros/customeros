import React, { useRef, useEffect, MouseEvent } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';
import { formatDistanceToNow } from 'date-fns/formatDistanceToNow';
import { differenceInCalendarMonths } from 'date-fns/differenceInCalendarMonths';

import { cn } from '@ui/utils/cn';
import { Select } from '@ui/form/Select';
import { Input } from '@ui/form/Input/Input';
import { Clock } from '@ui/media/icons/Clock';
import { Check } from '@ui/media/icons/Check';
import { File02 } from '@ui/media/icons/File02';
import { Mail01 } from '@ui/media/icons/Mail01';
import { User03 } from '@ui/media/icons/User03';
import { Social, Contact } from '@graphql/types';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { Share07 } from '@ui/media/icons/Share07';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useStore } from '@shared/hooks/useStore';
import { Calendar } from '@ui/media/icons/Calendar';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { SelectOption } from '@shared/types/SelectOptions';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { PhoneOutgoing02 } from '@ui/media/icons/PhoneOutgoing02';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { AutoresizeTextarea } from '@ui/form/Textarea/AutoresizeTextarea';
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
      store.contacts.remove(id);
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

    return (
      <>
        <Card
          ref={cardRef}
          key={contactStore?.id}
          className={cn(
            'bg-white w-full group rounded-lg border-[1px] border-gray-200 cursor-pointer hover:shadow-md ',
            isExpanded ? 'shadow-md' : 'shadow-xs',
            'ease-linear',
            'transition-all',
            'duration-1000',
          )}
        >
          <CardHeader onClick={toggle} className={cn('flex p-4 relative')}>
            <Avatar
              name={contactStore?.value.name ?? ''}
              src={
                contactStore?.value?.profilePhotoUrl
                  ? contactStore.value.profilePhotoUrl
                  : undefined
              }
              icon={<User03 className='text-primary-700 size-6' />}
              variant='shadowed'
            />

            <div className='ml-4 flex flex-col flex-1'>
              <Input
                className='font-semibold text-gray-700'
                name='name'
                size='xs'
                ref={nameInputRef}
                placeholder='Name'
                value={contactStore?.value?.name ?? ''}
                onChange={handleChange}
              />
              <Input
                className='text-gray-500'
                name='prefix'
                size='xs'
                placeholder='Title'
                value={contactStore?.value?.prefix ?? ''}
                onChange={handleChange}
              />
              <Select
                isMulti
                size='xs'
                name='role'
                options={roleOptions}
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
            {isExpanded && (
              <IconButton
                className='absolute z-50 top-2 right-2 p-1 opacity-0 pointer-events-auto transition-opacity duration-300 group-hover:opacity-100 "'
                size='xs'
                variant='ghost'
                colorScheme='gray'
                id='collapse-button'
                aria-label='Close'
                onClick={onClose}
                icon={<Check className='text-gray-500' />}
              />
            )}

            {!isExpanded && (
              <IconButton
                className='hover:bg-error-100 *:hover:text-error-500 absolute z-50 top-2 right-2 p-1 opacity-0 pointer-events-auto transition-opacity duration-300 group-hover:opacity-100 "'
                size='sm'
                variant='ghost'
                colorScheme='gray'
                id='confirm-button'
                aria-label='Delete contact'
                onClick={toggleConfirmDelete}
                icon={<Trash01 className='text-gray-400' />}
              />
            )}
          </CardHeader>
          {isExpanded && (
            <CardContent
              className={cn('flex flex-col', isExpanded ? 'h-auto' : 'h-0')}
            >
              <InputGroup>
                <LeftElement>
                  <Tooltip label='Click to autopopulate' hasArrow>
                    <span>
                      {contactStore?.isLoading ? (
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
                  ref={emailInputRef}
                  placeholder='Email'
                  variant='unstyled'
                  value={contactStore?.value?.emails?.[0]?.email ?? ''}
                  onChange={(e) => {
                    contactStore?.update(
                      (value) => {
                        set(value, 'emails[0].email', e.target.value);

                        return value;
                      },
                      { mutate: false },
                    );
                  }}
                  onBlur={() => {
                    if (!contactStore?.value?.emails?.[0]?.id) {
                      contactStore?.addEmail();
                    } else {
                      contactStore?.updateEmail();
                    }
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
                  value={
                    contactStore?.value.phoneNumbers?.[0]?.rawPhoneNumber ?? ''
                  }
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
                  onBlur={() => {
                    if (!contactStore?.value.phoneNumbers?.[0]?.id) {
                      contactStore?.addPhoneNumber();
                    } else {
                      contactStore?.updatePhoneNumber();
                    }
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

              <SocialIconInput
                placeholder='Social link'
                leftElement={<Share07 className='text-gray-500' />}
                value={
                  contactStore?.value?.socials?.map((s) => ({
                    label: s.url,
                    value: s.id,
                  })) ?? []
                }
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
                onCreate={(value) => {
                  contactStore?.update((prev) => {
                    prev.socials.push({
                      id: crypto.randomUUID(),
                      url: value,
                    } as Social);

                    return prev;
                  });
                }}
              />
              <TimezoneSelect
                isClearable
                placeholder='Timezone'
                options={timezoneOptions}
                value={timezoneOptions.find(
                  (v) => v.value === contactStore?.value?.timezone,
                )}
                onChange={(opt) => {
                  contactStore?.update((value) => {
                    value.timezone = opt?.value;

                    return value;
                  });
                }}
                leftElement={<Clock className='text-gray-500 mr-3' />}
              />
              <AutoresizeTextarea
                className='items-start'
                name='description'
                placeholder='Notes'
                onChange={handleChange}
                value={contactStore?.value?.description ?? ''}
                leftElement={<File02 className='text-gray-500 mt-1 mr-1' />}
              />
            </CardContent>
          )}
        </Card>
        <ConfirmDeleteDialog
          label='Delete this contact?'
          confirmButtonLabel='Delete contact'
          isOpen={isOpen}
          onClose={onClose}
          onConfirm={handleDelete}
          hideCloseButton
        />
      </>
    );
  },
);
