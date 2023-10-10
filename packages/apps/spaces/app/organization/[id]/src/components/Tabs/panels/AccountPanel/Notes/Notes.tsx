import { useState, useRef } from 'react';
import { useForm } from 'react-inverted-form';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { Card, CardBody, CardFooter } from '@ui/layout/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { OrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';

import { invalidateAccountDetailsQuery } from '../utils';
import { NotesDTO } from './Notes.dto';

interface NotesProps {
  id: string;
  data?: OrganizationAccountDetailsQuery['organization'] | null;
}

export const Notes = ({ data, id }: NotesProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isFocused, setIsFocused] = useState(false);
  const queryClient = useQueryClient();
  const client = getGraphQLClient();

  const updateOrganization = useUpdateOrganizationMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });

  const note = data?.note ?? '';

  useForm({
    formId: 'account-notes-form',
    defaultValues: {
      notes: note,
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        setIsFocused(false);
        updateOrganization.mutate({
          input: NotesDTO.toPayload({ id, note: action.payload.value }),
        });
      }
      return next;
    },
  });

  return (
    <Card
      p='4'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      boxShadow={isFocused ? 'md' : 'xs'}
      _hover={{
        boxShadow: 'md',
      }}
      transition='all 0.2s ease-out'
    >
      <CardBody as={Flex} p='0' w='full' align='center'>
        <FeaturedIcon>
          <Icons.File2 />
        </FeaturedIcon>
        <Heading ml='5' size='sm' color='gray.700'>
          Notes
        </Heading>
      </CardBody>
      <CardFooter as={Flex} flexDir='column' padding={0}>
        <Divider color='gray.200' my='4' />
        <FormAutoresizeTextarea
          name='notes'
          formId='account-notes-form'
          placeholder='Write some notes or anything related to this customer'
          spellCheck={false}
          onFocus={() => setIsFocused(true)}
        />
      </CardFooter>
    </Card>
  );
};
