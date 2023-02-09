import React, {
  ButtonHTMLAttributes,
  ChangeEvent,
  EventHandler,
  FC,
  FormEventHandler,
  useRef,
} from 'react';
import { Control, Controller, FieldValues } from 'react-hook-form';
import { Editor } from 'primereact/editor';
import { Button } from '../../atoms';
import { RichTextHeader } from '../rich-text-header';
import axios from 'axios';
import { toast } from 'react-toastify';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}
interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  fieldName: string;
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
  onCancel = () => null,
  onGetFieldValue,
}) => {
  const isEditMode = mode === NoteEditorModes.EDIT;

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) {
      return;
    }

    const formData = new FormData();
    formData.append('file', e.target.files[0]);
    axios
      .post(`/fs/file`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      .then((r: any) => {
        fetch(`/fs/file/${r.data.id}/download`)
          .then(async (response: any) => {
            const blob = await response.blob();

            const reader = new FileReader();
            reader.onload = function () {
              const dataUrl = reader.result as any;

              if (dataUrl) {
                onTextChange(
                  onGetFieldValue('htmlEnhanced') +
                    `<img width="400" src='${dataUrl}' alt='${r.data.id}'>`,
                );
              } else {
                toast.error(
                  'There was a problem on our side and we are doing our best to solve it!',
                );
              }
            };
            reader.readAsDataURL(blob);
          })
          .catch((reason: any) => {
            toast.error(
              'There was a problem on our side and we are doing our best to solve it!',
            );
          });
      })
      .catch((reason: any) => {
        toast.error(
          'There was a problem on our side and we are doing our best to solve it!',
        );
      });
  };

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: isEditMode ? '-17px -24px' : 0,
      }}
    >
      <Controller
        name='htmlEnhanced'
        control={fieldController}
        render={({ field }) => (
          <Editor
            style={{
              height: isEditMode ? 'auto' : '120px',
              borderColor: isEditMode && 'transparent',
            }}
            className='w-full h-full'
            headerTemplate={<RichTextHeader onFileChange={handleFileChange} />}
            value={field.value}
            onTextChange={onTextChange}
          />
        )}
      />
      <div
        className={` flex justify-content-end  ${
          isEditMode ? 'mb-3 mr-3' : 'mt-3'
        }`}
      >
        {isEditMode ? (
          <>
            <Button
              onClick={onCancel}
              className={`${isEditMode ? 'mb-3 mr-3' : ''}`}
            >
              Cancel
            </Button>
            <Button
              onClick={onSave}
              mode='primary'
              className={`${isEditMode ? 'mb-3 mr-3' : ''}`}
            >
              Save
            </Button>
          </>
        ) : (
          <>
            <Button
              onClick={onSave}
              mode='primary'
              className={`${isEditMode ? 'mb-3 mr-3' : ''}`}
            >
              Add note
            </Button>
          </>
        )}
      </div>
    </div>
  );
};
