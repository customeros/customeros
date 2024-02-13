import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { useRemirror } from '@remirror/react';
import { htmlToProsemirrorNode } from 'remirror';
import { useDebounce, useWillUnmount } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/layout/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { OrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';

import { NotesDTO } from './Notes.dto';
import { invalidateAccountDetailsQuery } from '../utils';

interface NotesProps {
  id: string;
  data?: OrganizationAccountDetailsQuery['organization'] | null;
}

export const Notes = ({ data, id }: NotesProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();
  const client = getGraphQLClient();
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const updateOrganization = useUpdateOrganizationMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });

  const updateNote = useDebounce((note) => {
    updateOrganization.mutate({
      input: NotesDTO.toPayload({ id, note }),
    });
  }, 800);
  const { setDefaultValues } = useForm({
    formId: 'account-notes-form',
    defaultValues: {
      notes: data?.note ?? '',
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        updateNote(action.payload.value);
      }
      if (action.type === 'FIELD_BLUR') {
        updateNote.flush();
      }

      return next;
    },
  });

  useEffect(() => {
    setDefaultValues({ notes: data?.note });
  }, [data?.note, data?.id]);

  useEffect(() => {
    if (data?.note) {
      const prosemirrorNodeValue = htmlToProsemirrorNode({
        schema: remirrorProps.state.schema,
        content: `${data?.note}`,
      });
      remirrorProps.getContext()?.setContent(prosemirrorNodeValue);
    }
  }, [data?.note]);

  useWillUnmount(() => {
    updateNote.flush();
  });

  return (
    <Card
      p='4'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      boxShadow={'xs'}
      _hover={{
        boxShadow: 'md',
      }}
      _focusWithin={{
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
        <Flex
          position='relative'
          sx={{
            '& .remirror-editor-wrapper': {
              height: '100%',
              minHeight: '100px',
            },
            '& .remirror-editor.ProseMirror': {
              minHeight: '100px',
            },
            '& a': {
              color: 'primary.600',
            },
            '& .test': {
              maxWidth: '200px',
            },
          }}
        >
          <RichTextEditor
            {...remirrorProps}
            placeholder='Write some notes or anything related to this customer'
            formId='account-notes-form'
            name='notes'
            showToolbar={false}
          ></RichTextEditor>
        </Flex>
      </CardFooter>
    </Card>
  );
};
