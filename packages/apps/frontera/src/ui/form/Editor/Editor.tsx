import type { EditorThemeClasses } from 'lexical';

import { useRef, useState, useEffect } from 'react';

import { $nodesOfType, LexicalEditor } from 'lexical';
import { $generateHtmlFromNodes } from '@lexical/html';
import { HistoryPlugin } from '@lexical/react/LexicalHistoryPlugin';
import { RichTextPlugin } from '@lexical/react/LexicalRichTextPlugin';
import { CheckListPlugin } from '@lexical/react/LexicalCheckListPlugin';
import { EditorRefPlugin } from '@lexical/react/LexicalEditorRefPlugin';
import { AutoFocusPlugin } from '@lexical/react/LexicalAutoFocusPlugin';
import { ContentEditable } from '@lexical/react/LexicalContentEditable';
import { LexicalErrorBoundary } from '@lexical/react/LexicalErrorBoundary';
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
  className?: string;
  placeholder?: string;
  defaultHtmlValue?: string;
  mentionsOptions?: string[];
  children?: React.ReactNode;
  hashtagsOptions?: SelectOption[];
  onChange?: (html: string) => void;
  onHashtagCreate?: (hashtag: string) => void;
  onHashtagSearch?: (q: string | null) => void;
  onHashtagsChange?: (hashtags: SelectOption[]) => void;
  editorRef?: React.MutableRefObject<LexicalEditor | null>;
}

export const Editor = ({
  children,
  onChange,
  editorRef,
  className,
  hashtagsOptions,
  mentionsOptions,
  onHashtagSearch,
  onHashtagCreate,
  onHashtagsChange,
  defaultHtmlValue,
  placeholder = 'Type something',
}: EditorProps) => {
  const editor = useRef<LexicalEditor>(null);
  const [floatingAnchorElem, setFloatingAnchorElem] =
    useState<HTMLDivElement>();
  const [isLinkEditMode, setIsLinkEditMode] = useState<boolean>(false);

  const initialConfig: InitialConfigType = {
    namespace: 'Timeline',
    theme,
    onError,
    nodes,
  };

  const onRef = (_floatingAnchorElem: HTMLDivElement) => {
    if (_floatingAnchorElem !== null) {
      setFloatingAnchorElem(_floatingAnchorElem);
    }
  };

  useEffect(() => {
    const dispose = editor?.current?.registerUpdateListener(
      ({ editorState }) => {
        editorRef = editor;

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
    <LexicalComposer initialConfig={initialConfig}>
      <EditorRefPlugin editorRef={editor} />
      <CheckListPlugin />
      <AutoLinkPlugin />
      <HistoryPlugin />
      <AutoFocusPlugin />
      <MentionsPlugin options={mentionsOptions} />
      <HashtagsPlugin
        onCreate={onHashtagCreate}
        onSearch={onHashtagSearch}
        options={hashtagsOptions ?? []}
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
              className='focus:outline-none'
              spellCheck='false'
            />
          </div>
        }
        placeholder={
          <span className='absolute top-0 text-gray-400'>{placeholder}</span>
        }
        ErrorBoundary={LexicalErrorBoundary}
      />
      {children}
    </LexicalComposer>
  );
};
