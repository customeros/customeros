'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import { useRef, useState, MouseEvent } from 'react';

import { useQueryClient } from '@tanstack/react-query';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import differenceInCalendarMonths from 'date-fns/differenceInCalendarMonths';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Contact } from '@graphql/types';
import { Avatar } from '@ui/media/Avatar';
import { useDisclosure } from '@ui/utils';
import { FormInput } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { Fade } from '@ui/transitions/Fade';
import { User01 } from '@ui/media/icons/User01';
import { IconButton } from '@ui/form/IconButton';
import { Collapse } from '@ui/transitions/Collapse';
import { FormInputGroup } from '@ui/form/InputGroup';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { SelectOption } from '@shared/types/SelectOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useUpdateContactMutation } from '@organization/src/graphql/updateContact.generated';
import { useDeleteContactMutation } from '@organization/src/graphql/deleteContact.generated';
import { useAddContactEmailMutation } from '@organization/src/graphql/addContactEmail.generated';
import { useAddContactSocialMutation } from '@organization/src/graphql/addContactSocial.generated';
import { useRemoveContactEmailMutation } from '@organization/src/graphql/removeContactEmail.generated';
import { useUpdateContactRoleMutation } from '@organization/src/graphql/updateContactJobRole.generated';
import { useAddContactPhoneNumberMutation } from '@organization/src/graphql/addContactPhoneNumber.generated';
import { useUpdateContactPhoneNumberMutation } from '@organization/src/graphql/updateContactPhoneNumber.generated';
import { useRemoveContactPhoneNumberMutation } from '@organization/src/graphql/removeContactPhoneNumber.generated';
import { EmailValidationMessage } from '@organization/src/components/Tabs/panels/PeoplePanel/ContactCard/EmailValidationMessage';

import { FormRoleSelect } from './FormRoleSelect';
import { FormTimezoneSelect } from './FormTimezoneSelect';
import { invalidateQuery, timezoneOptions } from '../util';
import { ContactForm, ContactFormDto } from './Contact.dto';
import { FormSocialInput } from '../../../shared/FormSocialInput';

interface ContactCardProps {
  contact: Contact;
  organizationName?: string;
}

export const ContactCard = ({
  contact,
  organizationName,
}: ContactCardProps) => {
  const client = getGraphQLClient();
  const organizationId = useParams()?.id as string;
  const queryClient = useQueryClient();
  const cardRef = useRef<HTMLDivElement>(null);
  const [isExpanded, setIsExpanded] = useState(false);
  const [roleIsFocused, setRoleIsFocused] = useState(false);
  const { isOpen, onOpen, onClose } = useDisclosure();

  useOutsideClick({
    ref: cardRef,
    handler: () => setIsExpanded(false),
  });

  const data = ContactFormDto.toForm(contact);

  const formId = `contact-form-${data.id}`;

  const invalidate = () => invalidateQuery(queryClient, organizationId);
  const updateContact = useUpdateContactMutation(client, {
    onSuccess: invalidate,
  });
  const updateRole = useUpdateContactRoleMutation(client, {
    onSuccess: invalidate,
  });
  const deleteContact = useDeleteContactMutation(client, {
    onSuccess: invalidate,
  });
  const addEmail = useAddContactEmailMutation(client, {
    onSuccess: invalidate,
  });
  const removeEmail = useRemoveContactEmailMutation(client, {
    onSuccess: invalidate,
  });
  const addPhoneNumber = useAddContactPhoneNumberMutation(client, {
    onSuccess: invalidate,
  });
  const updatePhoneNumber = useUpdateContactPhoneNumberMutation(client, {
    onSuccess: invalidate,
  });
  const removePhoneNumber = useRemoveContactPhoneNumberMutation(client, {
    onSuccess: invalidate,
  });
  const addSocial = useAddContactSocialMutation(client, {
    onSuccess: invalidate,
  });

  const toggle = (e: MouseEvent<HTMLDivElement>) => {
    if (['name', 'role', 'title'].includes((e.target as any)?.id)) {
      setIsExpanded(true);

      return;
    }
    setIsExpanded((prev) => !prev);
  };

  const prevEmail = data.email;
  const prevPhoneNumberId = data.phoneId;

  const timeAt = (() => {
    if (!data.startedAt) return undefined;
    const months = Math.abs(
      differenceInCalendarMonths(new Date(data.startedAt), new Date()),
    );

    if (months < 0) return `Less than a month at ${organizationName}`;
    if (months === 1) return `${months} month at ${organizationName}`;
    if (months > 1 && months < 12)
      return `${months} months at ${organizationName}`;
    if (months === 12) return `1 year at ${organizationName}`;
    if (months > 12)
      return `${formatDistanceToNow(
        new Date(data?.startedAt),
      )} at ${organizationName}`;
  })();

  const { state } = useForm<ContactForm>({
    formId,
    defaultValues: data,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        switch (action.payload.name) {
          case 'name':
          case 'timezone':
          case 'note': {
            updateContact.mutate(ContactFormDto.toDto({ ...state.values }));
            break;
          }
          case 'title':
          case 'role': {
            const key = (() => {
              const { name } = action.payload;
              if (name === 'role') return 'description';
              if (name === 'title') return 'jobTitle';

              return name;
            })();

            const value = (() => {
              if (action.payload.name === 'role') {
                return (action.payload.value as SelectOption[])
                  .map((v) => v.value)
                  .join(',');
              }

              return action.payload.value;
            })();

            updateRole.mutate({
              contactId: state.values.id,
              input: {
                id: state.values.roleId,
                description: state.values.role.map((v) => v.value).join(','),
                jobTitle: state.values.title,
                [key]: value,
              },
            });
            break;
          }
          case 'email': {
            const newEmail = action.payload.value;
            if (!newEmail) {
              removeEmail.mutate({ contactId: data.id, email: prevEmail });
              break;
            }
            addEmail.mutate({ contactId: data.id, input: { email: newEmail } });
            break;
          }
          case 'phone': {
            const newPhoneNumber = action.payload.value;
            if (!newPhoneNumber) {
              removePhoneNumber.mutate({
                contactId: data.id,
                id: prevPhoneNumberId,
              });
              break;
            }
            if (!prevPhoneNumberId) {
              addPhoneNumber.mutate({
                contactId: data.id,
                input: { phoneNumber: newPhoneNumber },
              });
              break;
            }
            updatePhoneNumber.mutate({
              contactId: data.id,
              input: {
                id: prevPhoneNumberId,
                phoneNumber: newPhoneNumber,
              },
            });
            break;
          }
          default:
            break;
        }
      }

      return next;
    },
  });

  const handleDelete = (e: MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();
    deleteContact.mutate({ contactId: data.id }, { onSuccess: onClose });
  };

  const toggleConfirmDelete = (e: MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();
    onOpen();
  };

  const handleAddSocial = ({
    newValue,
    onSuccess,
  }: {
    newValue: string;
    onSuccess: ({ id, url }: { id: string; url: string }) => void;
  }) => {
    addSocial.mutate(
      { contactId: contact.id, input: { url: newValue } },
      {
        onSuccess: ({ contact_AddSocial: { id, url } }) => {
          onSuccess({ id, url });
        },
      },
    );
  };

  return (
    <>
      <Card
        key={data.id}
        w='full'
        ref={cardRef}
        boxShadow={isExpanded ? 'md' : 'xs'}
        cursor='pointer'
        borderRadius='lg'
        border='1px solid'
        borderColor='gray.200'
        _hover={{
          boxShadow: 'md',
          '& > div > #confirm-button': {
            opacity: '1',
            pointerEvents: 'auto',
          },
        }}
        transition='all 0.2s ease-out'
      >
        <CardHeader
          as={Flex}
          p='4'
          pb={isExpanded ? 2 : 4}
          position='relative'
          onClick={toggle}
        >
          <Avatar
            name={state?.values?.name ?? data?.name}
            variant='shadowed'
            src={contact?.profilePhotoUrl ? contact.profilePhotoUrl : undefined}
            icon={<User01 color='gray.700' height='1.8rem' />}
          />
          <Flex ml='4' flexDir='column' flex='1'>
            <FormInput
              h='6'
              name='name'
              formId={formId}
              placeholder='Name'
              color='gray.700'
              fontWeight='semibold'
            />
            <FormInput
              h='6'
              name='title'
              color='gray.500'
              formId={formId}
              placeholder='Title'
            />
            <FormRoleSelect
              name='role'
              formId={formId}
              placeholder='Role'
              isCardOpen={isExpanded}
              isFocused={roleIsFocused}
              setIsFocused={setRoleIsFocused}
              data={data.role}
            />
          </Flex>
          {isExpanded && (
            <IconButton
              size='xs'
              top='2'
              right='2'
              variant='ghost'
              colorScheme='gray'
              id='collapse-button'
              position='absolute'
              aria-label='Close'
              onClick={onClose}
              icon={<Icons.Check color='gray.400' boxSize='5' />}
            />
          )}

          {!isExpanded && (
            <IconButton
              size='xs'
              top='2'
              right='2'
              variant='ghost'
              color='gray.400'
              colorScheme='gray'
              _hover={{
                background: 'red.100',
                color: 'red.400',
              }}
              opacity={0}
              pointerEvents='none'
              id='confirm-button'
              position='absolute'
              aria-label='Delete contact'
              isLoading={deleteContact.isLoading}
              onClick={toggleConfirmDelete}
              icon={<Icons.Trash1 boxSize='5' />}
            />
          )}
        </CardHeader>

        <Collapse
          in={isExpanded}
          style={{ overflow: 'unset' }}
          delay={{
            exit: 2,
          }}
        >
          <Fade
            in={isExpanded}
            delay={{
              enter: 0.2,
            }}
          >
            <CardBody pt={0}>
              <FormInputGroup
                formId={formId}
                name='email'
                placeholder='Email'
                leftElement={<Icons.Mail1 color='gray.500' />}
                rightElement={
                  <EmailValidationMessage
                    email={data.email}
                    validationDetails={
                      contact?.emails?.[0]?.emailValidationDetails
                    }
                  />
                }
              />
              <FormInputGroup
                formId={formId}
                name='phone'
                placeholder='Phone number'
                leftElement={<Icons.Phone2 color='gray.500' />}
              />
              {/* TODO: replace with FormInput. currently displayed as a text just for demoing purposes */}
              {timeAt && (
                <Flex align='center' h='39px'>
                  <Icons.Calendar color='gray.500' />
                  <Text ml='14px' cursor='text' textTransform='capitalize'>
                    {timeAt}
                  </Text>
                </Flex>
              )}
              {/* END TODO */}
              <FormSocialInput
                invalidateQuery={invalidate}
                addSocial={handleAddSocial}
                name='socials'
                formId={formId}
                placeholder='Social link'
                defaultValues={data?.socials}
                organizationId={organizationId}
                leftElement={<Icons.Share7 color='gray.500' />}
              />
              <FormTimezoneSelect
                formId={formId}
                isClearable
                name='timezone'
                placeholder='Timezone'
                options={timezoneOptions}
                leftElement={<Icons.Clock color='gray.500' mr='3' />}
              />
              <FormAutoresizeTextarea
                pl='30px'
                formId={formId}
                name='note'
                placeholder='Notes'
                leftElement={<Icons.File2 color='gray.500' />}
              />
            </CardBody>
          </Fade>
        </Collapse>
      </Card>
      <ConfirmDeleteDialog
        label='Delete this contact?'
        confirmButtonLabel='Delete contact'
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={handleDelete}
        isLoading={deleteContact.isLoading}
      />
    </>
  );
};
