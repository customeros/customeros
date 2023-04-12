import React, { FC, PropsWithChildren } from 'react';
import { IdentifierSchemaAttributes, prosemirrorNodeToHtml } from 'remirror';
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
import { SaveButtonWithOptions } from '../../atoms/button';
import { useFileData } from '../../../../hooks/useFileData';

import { NoteEditorModes } from './Editor';
import classNames from 'classnames';
import {
  UploadImageButton,
  Mention,
  CancelButton,
  CustomEditorToolbar,
} from './components';

export const extraAttributes: IdentifierSchemaAttributes[] = [
  {
    identifiers: ['mention', 'emoji'],
    attributes: { role: { default: 'presentation' } },
  },
  { identifiers: ['mention'], attributes: { href: { default: `/` } } },
];

export const SocialEditor: FC<PropsWithChildren<any>> = ({
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
  value,
  manager,
  state,
  setState,
  items,
  mode = 'ADD',
  editable = true,
  ...rest
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleAddFileToTextContent = (imagePreview: string) => {
    const data = prosemirrorNodeToHtml(state.doc);

    onHtmlChanged(data + imagePreview);
  };
  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  return (
    <div
      className={classNames(styles.editorWrapper, {
        [styles.editorWrapper]: !isEditMode,
        [styles.readOnly]: !editable,
        'remirror-read-only': !editable,
      })}
    >
      <Remirror
        editable={editable}
        manager={manager}
        state={state}
        onChange={(parameter) => {
          // Update the state to the latest value.
          setState(parameter.state);
        }}
      >
        <CustomEditorToolbar editable={editable} />
        <EditorComponent />
        <EmojiPopupComponent />
        <Mention />
        <TableComponents />

        <div
          className={classNames(styles.toolbar, {
            [styles.hidden]: !editable,
          })}
        >
          {children}
          <Toolbar>
            <div className={styles.toolbarActionButtons}>
              {!isEditMode && <CancelButton />}
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

            {!isEditMode && (
              <div className={styles.saveButtons}>
                <SaveButtonWithOptions
                  mode='primary'
                  items={items}
                  loading={saving}
                />
              </div>
            )}
          </Toolbar>
        </div>
      </Remirror>
    </div>
  );
};
