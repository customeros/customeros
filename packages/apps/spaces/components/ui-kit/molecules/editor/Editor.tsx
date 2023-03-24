import React, { ButtonHTMLAttributes, FC, ReactNode } from 'react';
import { Editor as PrimereactEditor } from 'primereact/editor';
import { RichTextHeader } from '../rich-text-header';
import { useFileData } from '../../../../hooks/useFileData';

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
  const handleAddFileToTextContent = (imagePreview: string) => {
    onHtmlChanged(onGetFieldValue('htmlEnhanced') + imagePreview);
  };

  const { onFileChange } = useFileData({
    addFileToTextContent: handleAddFileToTextContent,
  });

  return (
    <>
      {children}
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
