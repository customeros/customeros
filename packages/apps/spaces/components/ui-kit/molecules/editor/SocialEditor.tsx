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
  useRemirrorContext,
} from '@remirror/react';
import styles from './editor.module.scss';
import { Button, SaveButtonWithOptions } from '../../atoms/button';
import { useFileData } from '../../../../hooks/useFileData';

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
  onCancel,
  onPhoneCallSave,
  saving,
  onSave,
  manager,
  state,
  setState,
  mode = 'ADD',
  editable = true,
  onSubmit,
  submitButtonLabel,
  items,
  context,

  ...rest
}) => {
  const isEditMode = mode === 'EDIT';
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
      className={classNames(styles.editorWrapper, rest?.className, {
        [styles.editorWrapper]: !isEditMode,
        [styles.readOnly]: !editable,
        'remirror-read-only': !editable,
        'remirror-editable': editable,
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
              <>
                {items?.length ? (
                  <div className={styles.saveButtons}>
                    <Button onClick={onSubmit}>{submitButtonLabel}</Button>
                  </div>
                ) : (
                  <Button
                    onClick={onSubmit}
                    mode='primary'
                    style={{
                      padding: `0 8px`,
                      height: 32,
                      marginRight: '4px',
                      borderRadius: 4,
                    }}
                    className={styles.toolbarButton}
                  >
                    {submitButtonLabel}
                  </Button>
                )}
              </>
            )}
          </Toolbar>
        </div>
      </Remirror>
    </div>
  );
};
