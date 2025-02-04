import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { EditOrganizationNotesUsecase } from '@domain/usecases/organization-details/edit-organization-notes.usecase.ts';

import { File02 } from '@ui/media/icons/File02';
import { Editor } from '@ui/form/Editor/Editor';
import { Divider } from '@ui/presentation/Divider/Divider';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

interface NotesProps {
  id: string;
}

export const Notes = observer(({ id }: NotesProps) => {
  const editNotesUsecase = useMemo(
    () => new EditOrganizationNotesUsecase(id),
    [id],
  );

  return (
    <Card className='bg-white p-4 w-full cursor-default hover:shadow-md focus-within:shadow-md transition-all duration-200 ease-out'>
      <CardContent className='flex p-0 w-full items-center'>
        <FeaturedIcon colorScheme='gray' className='mr-4 ml-3 my-1 mt-3'>
          <File02 />
        </FeaturedIcon>
        <h2 className='ml-5 text-gray-700 font-semibold'>Notes</h2>
      </CardContent>
      <CardFooter className='flex flex-col items-start p-0 w-full'>
        <Divider className='my-4' />

        <div className='w-full'>
          <Editor
            size='sm'
            className='cursor-text'
            namespace='opportunity-next-step'
            placeholderClassName='cursor-text'
            onBlur={() => editNotesUsecase.execute()}
            dataTest='organization-account-notes-editor'
            defaultHtmlValue={editNotesUsecase.defaultNote ?? ''}
            placeholder='Write some notes or anything related to this customer'
            onChange={(html) => {
              editNotesUsecase.setNotes(html);
            }}
          />
        </div>
      </CardFooter>
    </Card>
  );
});
