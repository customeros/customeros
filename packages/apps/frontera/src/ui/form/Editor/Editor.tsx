import type { EditorThemeClasses } from 'lexical';

import React, {
  useRef,
  useState,
  useEffect,
  forwardRef,
  useImperativeHandle,
} from 'react';

import { HistoryPlugin } from '@lexical/react/LexicalHistoryPlugin';
import { $insertNodes, $nodesOfType, LexicalEditor } from 'lexical';
import { RichTextPlugin } from '@lexical/react/LexicalRichTextPlugin';
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

interface EditorProps {
  namespace: string;
  className?: string;
  placeholder?: string;
  defaultHtmlValue?: string;
  mentionsOptions?: string[];
  children?: React.ReactNode;
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
      onBlur,
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
      <div className='relative h-full'>
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
          <RichTextPlugin
            contentEditable={
              <div className={cn('relative', className)} ref={onRef}>
                <ContentEditable
                  onBlur={onBlur}
                  className='focus:outline-none'
                  spellCheck='false'
                />
              </div>
            }
            placeholder={
              <span
                className='absolute top-0 text-gray-400'
                onClick={() => editor.current?.focus()}
              >
                {placeholder}
              </span>
            }
            ErrorBoundary={LexicalErrorBoundary}
          />
          {children}
        </LexicalComposer>
      </div>
    );
  },
);
