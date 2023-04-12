'use client';
import React, { ButtonHTMLAttributes, FC, ReactNode, useCallback } from 'react';
import { Editor as PrimereactEditor } from 'primereact/editor';
import { RichTextHeader } from '../rich-text-header';
import { useFileData } from '../../../../hooks/useFileData';
import {
  BlockquoteExtension,
  BoldExtension,
  ImageExtension,
  ItalicExtension,
  LinkExtension,
  TextColorExtension,
  UnderlineExtension,
  FontSizeExtension,
  HistoryExtension,
  AnnotationExtension,
} from 'remirror/extensions';
import {
  useRemirror,
  useRemirrorContext,
  useHelpers,
  PlaceholderExtension,
} from '@remirror/react';
import { SocialEditor } from './SocialEditor';
import { Button, Send } from '../../atoms';
import styles from './editor.module.scss';
import { SaveButtonWithOptions } from '../../atoms/button';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  value: string;
  onGetFieldValue: (data: string) => string;
  mode: NoteEditorModes;
  onHtmlChanged: (html: string) => void;
  onSave: () => void;
  onPhoneCallSave?: () => void;
  onCancel?: () => void;
  label: string;
  possibleActions?: Array<{
    label: string;
    action: () => void;
  }>;
  children?: ReactNode;
  saving?: boolean;
}

const ALL_USERS = [
  { id: 'joe', label: 'Joe' },
  { id: 'sue', label: 'Sue' },
  { id: 'pat', label: 'Pat' },
  { id: 'tom', label: 'Tom' },
  { id: 'jim', label: 'Jim' },
];

const TAGS = ['editor', 'remirror', 'opensource', 'prosemirror'];

const SAMPLE_DOC = {
  type: 'doc',
  content: [
    {
      type: 'paragraph',
      attrs: { dir: null, ignoreBidiAutoUpdate: null },
      content: [{ type: 'text', text: 'Loaded content' }],
    },
  ],
};

function LoadButton() {
  const { setContent } = useRemirrorContext();
  const handleClick = useCallback(() => setContent(SAMPLE_DOC), [setContent]);

  return (
    <div>
      <Button
        className={styles.toolbarButton}
        mode='primary'
        onMouseDown={(event) => event.preventDefault()}
        onClick={handleClick}
      >
        Log call
      </Button>
    </div>
  );
}

function SaveButton() {
  const { getJSON } = useHelpers();
  const handleClick = useCallback(
    () => alert(JSON.stringify(getJSON())),
    [getJSON],
  );

  return (
    <div>
      <Button
        className={styles.toolbarButton}
        mode='primary'
        onMouseDown={(event) => event.preventDefault()}
        onClick={handleClick}
      >
        Log note
      </Button>
    </div>
  );
}

export const Editor: FC<Props> = ({
  mode,
  onHtmlChanged,
  onSave,
  onPhoneCallSave,
  label,
  value,
  onGetFieldValue,
  children,
  onCancel,
  saving = false,
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;
  const { manager, state } = useRemirror({
    extensions: () => [
      new BoldExtension(),
      new ItalicExtension(),
      new ImageExtension(),
      new LinkExtension(),
      new TextColorExtension(),
      new UnderlineExtension(),
      new FontSizeExtension(),
      new HistoryExtension(),
      new AnnotationExtension(),
      new BlockquoteExtension(),

      // new EmojiExtension({ data, plainText: true }),
      new PlaceholderExtension({ placeholder: `Type : to insert emojis` }),
    ],
    content: '<p>I love <b>Remirror</b></p>',
    selection: 'start',
    stringHandler: 'html',
  });

  const handleAddFileToTextContent = (imagePreview: string) => {
    onHtmlChanged(onGetFieldValue('htmlEnhanced') + imagePreview);
  };

  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  const items = [
    {
      label: 'Log as note',
      command: onPhoneCallSave,
    },
    {
      label: 'Log as phone call',
      command: () => {
        onHtmlChanged('');
      },
    },
    {
      label: '',
      command: onCancel,
    },
  ];

  return (
    <>
      <PrimereactEditor
        style={{
          height: isEditMode ? 'auto' : '160px',
          borderBottomColor: isEditMode && 'transparent',
        }}
        headerTemplate={
          <RichTextHeader
            hideButtons={isEditMode}
            onFileChange={(e) => onFileChange(e)}
            onSubmit={onSave}
            onCancel={onCancel}
            label={label}
            saving={saving}
            onSavePhoneCall={onPhoneCallSave}
          />
        }
        value={value}
        onTextChange={(e: any) => {
          onHtmlChanged(e.htmlValue);
        }}
      />
    </>
  );
};
