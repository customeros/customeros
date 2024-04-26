'use client';

import React from 'react';
import { useField } from 'react-inverted-form';

import { twMerge } from 'tailwind-merge';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import { useEditor, EditorContent, EditorContentProps } from '@tiptap/react';

interface RichTextEditorProps extends Omit<EditorContentProps, 'editor'> {
  name: string;
  formId: string;
  className?: string;
  placeholder: string;
}

export const RichTextEditor = ({
  placeholder,
  className,
  name,
  formId,
  ...props
}: RichTextEditorProps) => {
  const { getInputProps } = useField(name, formId);
  const { onChange, value } = getInputProps();

  const editor = useEditor({
    onUpdate: ({ editor }) => {
      const newValue = editor?.getHTML();
      onChange(newValue);
    },
    content: value,
    extensions: [
      StarterKit,
      Placeholder.configure({
        placeholder: placeholder,
      }),
    ],
  });

  return (
    <EditorContent
      className={twMerge('w-full h-full focus:outline-none', className)}
      editor={editor}
    />
  );
};
