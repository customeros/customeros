import type { EditorThemeClasses } from 'lexical';

import React, { useRef, useEffect } from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';
import { $getRoot, EditorState, LexicalEditor } from 'lexical';
import { HistoryPlugin } from '@lexical/react/LexicalHistoryPlugin';
import { OnChangePlugin } from '@lexical/react/LexicalOnChangePlugin';
import { EditorRefPlugin } from '@lexical/react/LexicalEditorRefPlugin';
import { AutoFocusPlugin } from '@lexical/react/LexicalAutoFocusPlugin';
import { ContentEditable } from '@lexical/react/LexicalContentEditable';
import { PlainTextPlugin } from '@lexical/react/LexicalPlainTextPlugin';
import { LexicalErrorBoundary } from '@lexical/react/LexicalErrorBoundary';
import {
  LexicalComposer,
  InitialConfigType,
} from '@lexical/react/LexicalComposer';

import { InputValidationPlugin } from '@ui/form/Editor/plugins/InputValidationPlugin.tsx';

const theme: EditorThemeClasses = {
  paragraph: 'my-0',
  text: {
    'invalid-email': 'invalid-email',
  },
};

const onError = (error: Error) => {
  console.error(error);
};

const contentEditableVariants = cva('focus:outline-none', {
  variants: {
    size: {
      xs: ['text-sm'],
      sm: ['text-sm'],
      md: ['text-base'],
      lg: ['text-lg'],
    },
  },
  defaultVariants: {
    size: 'md',
  },
});

// Extend the ParagraphNode to apply the custom format

interface EditorProps extends VariantProps<typeof contentEditableVariants> {
  namespace: string;
  dataTest?: string;
  className?: string;
  children?: React.ReactNode;
  type: 'linkedin' | 'email';
  placeholderClassName?: string;
  onChange?: (html: string) => void;
  onBlur?: (e: React.FocusEvent<HTMLDivElement>) => void;
  onKeyDown?: (e: React.KeyboardEvent<HTMLDivElement>) => void;
}

export const BulkContactsEditor = ({
  size,
  onBlur,
  dataTest,
  className,
  namespace,
  placeholderClassName,
  onKeyDown,
  type,
  onChange,
}: EditorProps) => {
  const editor = useRef<LexicalEditor | null>(null);

  const initialConfig: InitialConfigType = {
    namespace,
    theme,
    onError,
  };

  useEffect(() => {
    if (editor.current) {
      editor.current?.update(() => {
        const root = $getRoot();

        root.clear(); // Clear all nodes and reset to an empty editor state
      });
    }
  }, [type]);

  const onChangeHandler = (editorState: EditorState) => {
    const editorStateTextString = editorState.read(() => {
      return $getRoot().getTextContent();
    });

    if (onChange) {
      onChange(editorStateTextString);
    }
  };

  return (
    <div className='relative w-full h-full lexical-editor cursor-text min-h-[88px]'>
      <LexicalComposer initialConfig={initialConfig}>
        <EditorRefPlugin editorRef={editor} />
        <HistoryPlugin />
        <AutoFocusPlugin />
        <OnChangePlugin ignoreSelectionChange onChange={onChangeHandler} />
        <InputValidationPlugin type={type} />

        <PlainTextPlugin
          ErrorBoundary={LexicalErrorBoundary}
          contentEditable={
            <ContentEditable
              onBlur={onBlur}
              spellCheck='false'
              data-test={dataTest}
              className={twMerge(contentEditableVariants({ size, className }))}
              onKeyDown={(e) =>
                onKeyDown ? onKeyDown(e) : e.stopPropagation()
              }
            />
          }
          placeholder={
            <div
              onClick={() => editor.current?.focus()}
              className={twMerge(
                contentEditableVariants({
                  size,
                  className: placeholderClassName,
                }),
                'absolute top-0 text-gray-400 p-2 text-sm',
              )}
            >
              {type === 'email' ? (
                <>
                  <p>ivy@green.tech</p>
                  <p>pete@moss.ai</p>
                  <p>sue@flay.com</p>
                </>
              ) : (
                <>
                  <p>linkedin.com/in/ivy-green</p>
                  <p>linkedin.com/in/pete-moss</p>
                  <p>linkedin.com/in/sue-flay</p>
                </>
              )}
            </div>
          }
        />
      </LexicalComposer>
    </div>
  );
};
