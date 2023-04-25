import React, { FC, PropsWithChildren, useCallback, useEffect } from 'react';
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
import { Check, IconButton, Pencil } from '../../atoms';
import { useDebouncedCallback } from 'use-debounce';
import { editorMode } from '../../../../state';
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
  onCancel,
  saving,
  onSave,
  items,
  context,
  onDebouncedSave,
  isEditMode = false,
  onToggleEditMode,
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

  useEffect(() => {
    return () => {
      debounced.flush();
    };
  }, []);

  // console.log('🏷️ ----- isEditMode: ', isEditMode);

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
          // const html = prosemirrorNodeToHtml(parameter.state.doc);
          // debounced(html);
        }}
      >
        <CustomEditorToolbar editable={isEditMode} />
        <EditorComponent />
        <EmojiPopupComponent />
        <Mention />
        <TableComponents />

        <div className={classNames(styles.toolbar, styles.debouncedToolbar)}>
          {children}
          <Toolbar>
            <div
              className={classNames(styles.toolbarActionButtons, {
                [styles.hidden]: !isEditMode,
              })}
            >
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
            <div />

            <div className={styles.saveButtons}>
              {isEditMode && (
                <IconButton
                  isSquare
                  size='xxs'
                  className={styles.toolbarButton}
                  mode='subtle'
                  style={{ background: 'transparent' }}
                  onClick={() => onToggleEditMode(false)}
                  icon={<Check color={'#29C76F'} />}
                />
              )}

              {!isEditMode && (
                <IconButton
                  isSquare
                  size='xxs'
                  className={styles.toolbarButton}
                  mode='subtle'
                  onClick={() => onToggleEditMode(true)}
                  icon={<Pencil style={{ transform: 'scale(0.8)' }} />}
                />
              )}
            </div>
          </Toolbar>
        </div>
      </Remirror>
    </div>
  );
};
