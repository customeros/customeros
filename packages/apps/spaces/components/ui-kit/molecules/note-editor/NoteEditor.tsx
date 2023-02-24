import React, { ButtonHTMLAttributes, FC } from 'react';
import { Editor } from 'primereact/editor';
import { RichTextHeader } from '../rich-text-header';
import styles from '../../../../pages/contact/contact.module.scss';
import { useFileData } from '../../../../hooks/useFileData';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}
interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  value: string;
  onGetFieldValue: (data: string, imagePreview: string) => string;
  mode: NoteEditorModes;
  onTextChange: (e: any) => void;
  onSave: () => void;
}

export const NoteEditor: FC<Props> = ({
  mode,
  onTextChange,
  onSave,
  value,
  onGetFieldValue,
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleAddFileToTextContent = (imagePreview: string) => {
    onTextChange(onGetFieldValue('htmlEnhanced', imagePreview));
  };

  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  return (
    <form className={styles.editor}>
      <Editor
        style={{
          height: isEditMode ? 'auto' : '120px',
          borderColor: isEditMode && 'transparent',
        }}
        headerTemplate={
          <RichTextHeader
            onFileChange={(e) => onFileChange(e)}
            onSubmit={onSave}
          />
        }
        value={value}
        onTextChange={onTextChange}
      />
    </form>
  );
};
