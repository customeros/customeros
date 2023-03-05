import React, { ButtonHTMLAttributes, FC, ReactNode } from 'react';
import { Editor as PrimereactEditor } from 'primereact/editor';
import { RichTextHeader } from '../rich-text-header';
import styles from './editor.module.scss';
import { useFileData } from '../../../../hooks/useFileData';
import { Button } from '../../atoms';

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
  onCancel?: () => void;
  label: string;
  possibleActions?: Array<{
    label: string;
    action: () => void;
  }>;
  children?: ReactNode;
}

export const Editor: FC<Props> = ({
  mode,
  onTextChange,
  onSave,
  label,
  value,
  onGetFieldValue,
  children,
  onCancel,
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;
  const handleAddFileToTextContent = (imagePreview: string) => {
    onTextChange(onGetFieldValue('htmlEnhanced', imagePreview));
  };

  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  return (
    <>
      {children}
      <PrimereactEditor
        style={{
          height: isEditMode ? 'auto' : '120px',
          borderColor: isEditMode && 'transparent',
        }}
        headerTemplate={
          <RichTextHeader
            onFileChange={(e) => onFileChange(e)}
            onSubmit={onSave}
            onCancel={onCancel}
            label={label}
          />
        }
        value={value}
        onTextChange={onTextChange}
      />
    </>
  );
};
