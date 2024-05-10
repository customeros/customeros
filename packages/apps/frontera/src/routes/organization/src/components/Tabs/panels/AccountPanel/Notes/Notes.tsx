import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { useRemirror } from '@remirror/react';
import { htmlToProsemirrorNode } from 'remirror';
import { useDebounce, useWillUnmount } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';

import { File02 } from '@ui/media/icons/File02';
import { Divider } from '@ui/presentation/Divider/Divider';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { RichTextEditor } from '@ui/form/RichTextEditor2/RichTextEditor';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';
import { useUpdateOrganizationMutation } from '@shared/graphql/updateOrganization.generated';
import { OrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';

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
  useForm({
    formId: 'account-notes-form',
    defaultValues: {
      notes: data?.note ?? '<p style=""></p>',
    },
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        updateNote(action.payload.value);
      }
      if (action.type === 'FIELD_BLUR') {
        updateNote.flush();
      }

      return next;
    },
  });

  useEffect(() => {}, [data?.note, data?.id]);

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
    <Card className='bg-white p-4 w-full cursor-default hover:shadow-md focus-within:shadow-md transition-all duration-200 ease-out'>
      <CardContent className='flex p-0 w-full items-center'>
        <FeaturedIcon className='mr-4 ml-3 my-1 mt-3' colorScheme='gray'>
          <File02 />
        </FeaturedIcon>
        <h2 className='ml-5 text-gray-700 font-semibold '>Notes</h2>
      </CardContent>
      <CardFooter className='flex flex-col items-start p-0 w-full'>
        <Divider className='my-4' />

        <RichTextEditor
          formId='account-notes-form'
          name='notes'
          placeholder='Write some notes or anything related to this customer'
          className='min-h-[100px] cursor-text'
        />
      </CardFooter>
    </Card>
  );
};
