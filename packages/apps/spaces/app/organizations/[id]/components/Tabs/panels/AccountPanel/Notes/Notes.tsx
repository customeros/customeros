import { useState, useRef } from 'react';
import { useForm } from 'react-inverted-form';
import { useQueryClient } from '@tanstack/react-query';

import { Note } from '@graphql/types';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { Card, CardBody, CardFooter } from '@ui/layout/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useAddOrganizationNoteMutation } from '@organization/graphql/addOrganizationNote.generated';
import { useUpdateOrganizationNoteMutation } from '@organization/graphql/updateOrganizationNote.generated';

import { invalidateAccountDetailsQuery } from '../utils';

interface NotesProps {
  id: string;
  data?: Pick<Note, '__typename' | 'id' | 'content' | 'contentType'>[];
}

export const Notes = ({ data, id }: NotesProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isFocused, setIsFocused] = useState(false);
  const queryClient = useQueryClient();
  const client = getGraphQLClient();

  const addNote = useAddOrganizationNoteMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });
  const updateNote = useUpdateOrganizationNoteMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });

  const note = data?.[0];

  useForm({
    formId: 'account-notes-form',
    defaultValues: {
      notes: note?.content || '',
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        setIsFocused(false);

        if (!note) {
          addNote.mutate({
            organzationId: id,
            input: {
              content: action.payload.value,
              contentType: 'text',
            },
          });
        } else {
          updateNote.mutate({
            input: {
              id: note.id,
              content: action.payload.value,
            },
          });
        }
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
