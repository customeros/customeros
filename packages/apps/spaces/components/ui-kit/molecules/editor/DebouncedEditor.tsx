import React, {
  FC,
  PropsWithChildren,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { IdentifierSchemaAttributes, prosemirrorNodeToHtml } from 'remirror';
import { TableExtension } from '@remirror/extension-react-tables';
import { useDebouncedCallback } from 'use-debounce';
import {
  EditorComponent,
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
import { useFileData } from '@spaces/hooks/useFileData';

import classNames from 'classnames';
import { UploadImageButton, Mention, CustomEditorToolbar } from './components';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
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
import { useRemirror } from '@remirror/react';
import Pencil from '@spaces/atoms/icons/Pencil';
import Check from '@spaces/atoms/icons/Check';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
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
    ...wysiwygPreset(),
  ];
  const extensions = useCallback(() => [...remirrorExtentions], []);
  const [isFocused, setIsFocused] = useState(false);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: value,
  });

  const debounced = useDebouncedCallback(
    (value) => {
      const html = prosemirrorNodeToHtml(value);
      onDebouncedSave(html);
    },
    // delay in ms
    600,
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
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        onChange={(parameter) => {
          // Update the state to the latest value.
          setState(parameter.state);
          if (isFocused) {
            debounced(parameter.state.doc);
          }
        }}
      >
        <CustomEditorToolbar editable={isEditMode} />
        <EditorComponent />
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
              <IconButton
                isSquare
                label='Save'
                size='xxs'
                className={styles.toolbarButton}
                mode='subtle'
                style={{ background: 'transparent' }}
                onClick={() => onToggleEditMode(!isEditMode)}
                icon={
                  isEditMode ? (
                    <Check height={20} width={20} color='green' />
                  ) : (
                    <Pencil height={20} width={20} />
                  )
                }
              />
            </div>
          </Toolbar>
        </div>
      </Remirror>
    </div>
  );
};
