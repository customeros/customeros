import type { EditorThemeClasses } from 'lexical';

import React, {
  useRef,
  useState,
  useEffect,
  forwardRef,
  useImperativeHandle,
} from 'react';

import { twMerge } from 'tailwind-merge';
import { cva, VariantProps } from 'class-variance-authority';
import { HistoryPlugin } from '@lexical/react/LexicalHistoryPlugin';
import { $insertNodes, $nodesOfType, LexicalEditor } from 'lexical';
import { RichTextPlugin } from '@lexical/react/LexicalRichTextPlugin';
import { PlainTextPlugin } from '@lexical/react/LexicalPlainTextPlugin';
import { CheckListPlugin } from '@lexical/react/LexicalCheckListPlugin';
import { EditorRefPlugin } from '@lexical/react/LexicalEditorRefPlugin';
import { AutoFocusPlugin } from '@lexical/react/LexicalAutoFocusPlugin';
import { ContentEditable } from '@lexical/react/LexicalContentEditable';
import { LexicalErrorBoundary } from '@lexical/react/LexicalErrorBoundary';
import { $generateNodesFromDOM, $generateHtmlFromNodes } from '@lexical/html';
import {
  LexicalComposer,
  InitialConfigType,
} from '@lexical/react/LexicalComposer';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';

import { nodes } from './nodes/nodes';
import { HashtagNode } from './nodes/HashtagNode';
import MentionsPlugin from './plugins/MentionsPlugin';
import AutoLinkPlugin from './plugins/AutoLinkPlugin';
import HashtagsPlugin from './plugins/HashtagsPlugin';
import FloatingLinkEditorPlugin from './plugins/FloatingLinkEditorPlugin';

const theme: EditorThemeClasses = {};

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

interface EditorProps extends VariantProps<typeof contentEditableVariants> {
  namespace: string;
  dataTest?: string;
  className?: string;
  placeholder?: string;
  usePlainText?: boolean;
  defaultHtmlValue?: string;
  mentionsOptions?: string[];
  children?: React.ReactNode;
  placeholderClassName?: string;
  hashtagsOptions?: SelectOption[];
  onChange?: (html: string) => void;
  onHashtagCreate?: (hashtag: string) => void;
  onHashtagSearch?: (q: string | null) => void;
  onMentionsSearch?: (q: string | null) => void;
  onHashtagsChange?: (hashtags: SelectOption[]) => void;
  onBlur?: (e: React.FocusEvent<HTMLDivElement>) => void;
}

export const Editor = forwardRef<LexicalEditor | null, EditorProps>(
  (
    {
      size,
      onBlur,
      dataTest,
      children,
      onChange,
      className,
      namespace,
      onHashtagSearch,
      onHashtagCreate,
      onHashtagsChange,
      onMentionsSearch,
      defaultHtmlValue,
      hashtagsOptions = [],
      mentionsOptions = [],
      usePlainText = false,
      placeholderClassName,
      placeholder = 'Type something',
    },
    ref,
  ) => {
    const editor = useRef<LexicalEditor | null>(null);
    const hasLoadedDefaultHtmlValue = useRef(false);
    const [floatingAnchorElem, setFloatingAnchorElem] =
      useState<HTMLDivElement>();
    const [isLinkEditMode, setIsLinkEditMode] = useState<boolean>(false);

    const initialConfig: InitialConfigType = {
      namespace,
      theme,
      onError,
      nodes,
    };

    const EditorPlugin = usePlainText ? PlainTextPlugin : RichTextPlugin;

    const onRef = (_floatingAnchorElem: HTMLDivElement) => {
      if (_floatingAnchorElem !== null) {
        setFloatingAnchorElem(_floatingAnchorElem);
      }
    };

    useImperativeHandle(ref, () => editor.current as LexicalEditor);

    useEffect(() => {
      editor.current?.update(() => {
        if (!editor?.current || hasLoadedDefaultHtmlValue.current) return;

        if (defaultHtmlValue) {
          const parser = new DOMParser();
          const dom = parser.parseFromString(defaultHtmlValue, 'text/html');
          const nodes = $generateNodesFromDOM(editor?.current, dom);

          $insertNodes(nodes);
          hasLoadedDefaultHtmlValue.current = true;
        }
      });

      const dispose = editor?.current?.registerUpdateListener(
        ({ editorState }) => {
          editorState.read(() => {
            if (!editor?.current) return;

            const hashtagNodes = $nodesOfType(HashtagNode);
            const html = $generateHtmlFromNodes(editor?.current);

            onChange?.(html);
            onHashtagsChange?.(hashtagNodes.map((node) => node.__hashtag));
          });
        },
      );

      return () => {
        dispose?.();
      };
    }, []);

    return (
      <div className='relative w-full h-full'>
        <LexicalComposer initialConfig={initialConfig}>
          <EditorRefPlugin editorRef={editor} />
          <CheckListPlugin />
          <AutoLinkPlugin />
          <HistoryPlugin />
          <AutoFocusPlugin />
          <MentionsPlugin
            options={mentionsOptions}
            onSearch={onMentionsSearch}
          />
          <HashtagsPlugin
            options={hashtagsOptions}
            onCreate={onHashtagCreate}
            onSearch={onHashtagSearch}
          />
          {floatingAnchorElem && (
            <FloatingLinkEditorPlugin
              anchorElem={floatingAnchorElem}
              isLinkEditMode={isLinkEditMode}
              setIsLinkEditMode={setIsLinkEditMode}
            />
          )}
          <EditorPlugin
            ErrorBoundary={LexicalErrorBoundary}
            placeholder={
              <span
                onClick={() => editor.current?.focus()}
                className={twMerge(
                  contentEditableVariants({
                    size,
                    className: placeholderClassName,
                  }),
                  'absolute top-0 text-gray-400',
                )}
              >
                {placeholder}
              </span>
            }
            contentEditable={
              <div ref={onRef} className={cn('relative', className)}>
                <ContentEditable
                  onBlur={onBlur}
                  spellCheck='false'
                  data-test={dataTest}
                  onKeyDown={(e) => e.stopPropagation()}
                  className={twMerge(
                    contentEditableVariants({ size, className }),
                  )}
                />
              </div>
            }
          />
          {children}
        </LexicalComposer>
      </div>
    );
  },
);
