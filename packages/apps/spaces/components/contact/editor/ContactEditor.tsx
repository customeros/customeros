import React, { FC } from 'react';
import { Editor } from '../../ui-kit/molecules';
import { Controller, useForm } from 'react-hook-form';
import { useCreateContactNote } from '../../../hooks/useNote';
import { useRecoilState, useRecoilValue, useSetRecoilState } from 'recoil';
import { editorEmail, editorMode, EditorMode, userData } from '../../../state';
import { EmailFields } from './email-fields';
import classNames from 'classnames';
import { useCreatePhoneCallInteractionEvent } from '../../../hooks/useContact/useCreatePhoneCallInteractionEvent';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}

interface Props {
  contactId: string;
}

const DEFAULT_VALUES = {
  html: '',
  htmlEnhanced: '',
};
export const ContactEditor: FC<Props> = ({ contactId }) => {
  const { handleSubmit, setValue, getValues, control, reset } = useForm({
    defaultValues: DEFAULT_VALUES,
  });
  const { identity: loggedInUserEmail } = useRecoilValue(userData);
  const [editorModeState, setMode] = useRecoilState(editorMode);
  const {
    handleSubmit: handleSendEmail,
    to,
    respondTo,
  } = useRecoilValue(editorEmail);
  const { onCreateContactNote, saving } = useCreateContactNote({ contactId });
  const { onCreatePhoneCallInteractionEvent } =
    useCreatePhoneCallInteractionEvent({ contactId });

  const onSubmit = handleSubmit(async (d) => {
    //remove src attribute to not send the file bytes in here
    // identiti header - from session uuid
    const dataToSubmit = {
      appSource: 'Openline',
      html: d?.htmlEnhanced?.replaceAll(/.src(\S*)/g, '') || '',
    };

    editorModeState.mode === EditorMode.Email && handleSendEmail
      ? handleSendEmail(
          dataToSubmit.html.replace(/(<([^>]+)>)/gi, ''),
          () => reset(DEFAULT_VALUES),
          to,
          respondTo,
        )
      : onCreateContactNote(dataToSubmit).then(() => reset(DEFAULT_VALUES));
  });

  const onPhoneCallSubmit = handleSubmit(async (d) => {
    //remove src attribute to not send the file bytes in here
    // identiti header - from session uuid
    const dataToSubmit = {
      appSource: 'Openline',
      sentBy: loggedInUserEmail,
      content: d?.htmlEnhanced?.replaceAll(/.src(\S*)/g, '') || '',
      contentType: 'text/html',
    };

    onCreatePhoneCallInteractionEvent(dataToSubmit).then(() =>
      reset(DEFAULT_VALUES),
    );
  });

  const handleCancel = () => {
    setMode({ mode: EditorMode.Note, submitButtonLabel: 'Log as note' });
    reset(DEFAULT_VALUES);
  };

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: 0,
      }}
      key={editorModeState.mode}
    >
      <Controller
        name='htmlEnhanced'
        control={control}
        render={({ field }) => (
          <div
            className={classNames({
              'openline-editor-email':
                editorModeState.mode === EditorMode.Email,
            })}
          >
            {editorModeState.mode === EditorMode.Email && <EmailFields />}
            <Editor
              mode={NoteEditorModes.ADD}
              onGetFieldValue={getValues}
              value={field.value}
              saving={saving}
              onSave={onSubmit}
              onPhoneCallSave={onPhoneCallSubmit}
              onCancel={
                editorModeState.mode === EditorMode.Email
                  ? handleCancel
                  : undefined
              }
              label={editorModeState.submitButtonLabel}
              onHtmlChanged={(html: string) => setValue('htmlEnhanced', html)}
            />
          </div>
        )}
      />
    </div>
  );
};
