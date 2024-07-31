import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import { useEditor, EditorContent } from '@tiptap/react';

import { File02 } from '@ui/media/icons/File02';
import { useStore } from '@shared/hooks/useStore';
import { Divider } from '@ui/presentation/Divider/Divider';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

interface NotesProps {
  id: string;
}

export const Notes = observer(({ id }: NotesProps) => {
  const store = useStore();
  const organization = store.organizations.value.get(id);

  const editor = useEditor({
    onUpdate: ({ editor }) => {
      const newValue = editor?.getHTML();

      organization?.update((org) => {
        org.notes = newValue;

        return org;
      });
    },
    content: organization?.value?.notes,
    extensions: [
      StarterKit,
      Placeholder.configure({
        placeholder: 'Write some notes or anything related to this customer',
      }),
    ],
  });

  useEffect(() => {
    editor?.commands.setContent(organization?.value?.notes || '', false, {
      preserveWhitespace: 'full',
    });
  }, [organization?.value?.notes, editor]);

  return (
    <Card className='bg-white p-4 w-full cursor-default hover:shadow-md focus-within:shadow-md transition-all duration-200 ease-out'>
      <CardContent className='flex p-0 w-full items-center'>
        <FeaturedIcon colorScheme='gray' className='mr-4 ml-3 my-1 mt-3'>
          <File02 />
        </FeaturedIcon>
        <h2 className='ml-5 text-gray-700 font-semibold '>Notes</h2>
      </CardContent>
      <CardFooter className='flex flex-col items-start p-0 w-full'>
        <Divider className='my-4' />
        <EditorContent
          size={100}
          editor={editor}
          className='min-h-[100px] cursor-text w-full'
        />
      </CardFooter>
    </Card>
  );
});
