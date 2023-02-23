import React, { ButtonHTMLAttributes, ChangeEvent, FC } from 'react';
import { Control, Controller, FieldValues } from 'react-hook-form';
import { Editor } from 'primereact/editor';
import { Button } from '../../atoms';
import { RichTextHeader } from '../rich-text-header';
import axios from 'axios';
import { toast } from 'react-toastify';
import styles from '../../../../pages/contact/contact.module.scss';
import { useFileData } from '../../../../hooks/useFileData';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}
interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  fieldName: string;
  value: string;
  fieldController: Control<{ id: any; html: any; htmlEnhanced: any }, any>;
  onGetFieldValue: any;
  mode: NoteEditorModes;
  onTextChange: (e: any) => void;
  onSave: () => void;
  onCancel?: () => void;
}

export const NoteEditor: FC<Props> = ({
  fieldName,
  fieldController,
  mode,
  onTextChange,
  onSave,
  value,
  onCancel = () => null,
  onGetFieldValue,
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleAddFileToTextContent = (imagePreview: string) => {
    onTextChange(onGetFieldValue('htmlEnhanced', imagePreview));
  };

  const { onFileChange } = useFileData(handleAddFileToTextContent);

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
