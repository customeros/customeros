import React, { FC, PropsWithChildren, useCallback } from 'react';
import { IdentifierSchemaAttributes, prosemirrorNodeToHtml } from 'remirror';
import { TableExtension } from '@remirror/extension-react-tables';

import {
  EditorComponent,
  EmojiPopupComponent,
  Remirror,
  TableComponents,
  ToggleBlockquoteButton,
  Toolbar,
  HistoryButtonGroup,
  HeadingLevelButtonGroup,
  CommandButtonGroup,
  ToggleBulletListButton,
  ToggleOrderedListButton,
  ToggleTaskListButton,
  CreateTableButton,
} from '@remirror/react';
import styles from './editor.module.scss';
import { useFileData } from '../../../../hooks/useFileData';

import classNames from 'classnames';
import { UploadImageButton, Mention, CustomEditorToolbar } from './components';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  EmojiExtension,
  FontSizeExtension,
  HistoryExtension,
  ImageExtension,
  ItalicExtension,
  LinkExtension,
  MentionAtomExtension,
  OrderedListExtension,
  StrikeExtension,
  TextColorExtension,
  UnderlineExtension,
  wysiwygPreset,
} from 'remirror/extensions';
import data from 'svgmoji/emoji.json';
import { useRemirror } from '@remirror/react';
import { Check, IconButton } from '../../atoms';
import { useDebouncedCallback } from 'use-debounce';
export const extraAttributes: IdentifierSchemaAttributes[] = [
  {
    identifiers: ['mention', 'emoji'],
    attributes: { role: { default: 'presentation' } },
  },
  { identifiers: ['mention'], attributes: { href: { default: `/` } } },
];

export const DebouncedEditor: FC<PropsWithChildren<any>> = ({
  placeholder,
  stringHandler,
  children,
  users,
  tags,
  onHtmlChanged,
  onPhoneCallSave,
  onCancel,
  saving,
  onSave,
  items,
  context,
  onDebouncedSave,
  isEditMode = false,
  className,
  value = '',
  ...rest
}) => {
  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

    new EmojiExtension({ plainText: true, data, moji: 'noto' }),
    ...wysiwygPreset(),
    new BoldExtension(),
    new ItalicExtension(),
    new BlockquoteExtension(),
    new ImageExtension({}),
    new LinkExtension({ autoLink: true }),
    new TextColorExtension(),
    new UnderlineExtension(),
    new FontSizeExtension(),
    new HistoryExtension(),
    new AnnotationExtension(),
    new BulletListExtension(),
    new OrderedListExtension(),
    new StrikeExtension(),
  ];
  const extensions = useCallback(() => [...remirrorExtentions], []);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: value,
  });
  const debounced = useDebouncedCallback(
    // function
    (value) => {
      onDebouncedSave(value);
    },
    // delay in ms
    300,
  );
  const handleAddFileToTextContent = (imagePreview: string) => {
    const data = prosemirrorNodeToHtml(state.doc);
    const htmlData = data + imagePreview;
    context.setContent(data + imagePreview);
  };
  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  return (
    <div
      className={classNames(
        styles.editorWrapper,
        styles.debouncedEditor,
        'remirror-debounced-editor',
        {
          [className]: !!className,
          'remirror-read-only': !isEditMode,
          'remirror-editable': isEditMode,
        },
      )}
    >
      <Remirror
        editable={isEditMode}
        manager={manager}
        state={state}
        onChange={(parameter) => {
          // Update the state to the latest value.
          setState(parameter.state);
          const html = prosemirrorNodeToHtml(parameter.state.doc);
          onDebouncedSave(html);
        }}
      >
        <CustomEditorToolbar editable={isEditMode} />
        <EditorComponent />
        <EmojiPopupComponent />
        <Mention />
        <TableComponents />

        <div
          className={classNames(styles.toolbar, {
            [styles.hidden]: !isEditMode,
          })}
        >
          {children}
          <Toolbar>
            <div className={styles.toolbarActionButtons}>
              <HistoryButtonGroup />
              <HeadingLevelButtonGroup />
              <ToggleBlockquoteButton />

              <CommandButtonGroup>
                <ToggleBulletListButton />
                <ToggleOrderedListButton />
                <ToggleTaskListButton />
                <CreateTableButton />
              </CommandButtonGroup>
              <UploadImageButton onFileChange={onFileChange} />
            </div>

            {/*{!isEditMode && (*/}
            <div className={styles.saveButtons}>
              <IconButton
                isSquare
                mode='success'
                size='xxs'
                onClick={() => console.log('')}
                icon={<Check />}
              />
            </div>
            {/*)}*/}
          </Toolbar>
        </div>
      </Remirror>
    </div>
  );
};
