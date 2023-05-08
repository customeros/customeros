import React, { FC, useCallback, useRef } from 'react';
import { useRecoilState, useRecoilValue } from 'recoil';
import { editorEmail, editorMode, EditorMode, userData } from '../../../state';
import { useCreatePhoneCallInteractionEvent } from '../../../hooks/useContact/useCreatePhoneCallInteractionEvent';
import {
  extraAttributes,
  SocialEditor,
} from '../../ui-kit/molecules/editor/SocialEditor';
import { prosemirrorNodeToHtml } from 'remirror';
import { useRemirror } from '@remirror/react';
import { TableExtension } from '@remirror/extension-react-tables';
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
import { toast } from 'react-toastify';
import { EmailFields } from './email-fields';

export enum NoteEditorModes {
  'ADD' = 'ADD',
  'EDIT' = 'EDIT',
}

interface Props {
  contactId: string;
}

export const ContactEditor: FC<Props> = ({ contactId }) => {
  const [editorModeState, setMode] = useRecoilState(editorMode);

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

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
  const extensions = useCallback(
    () => [...remirrorExtentions],
    [editorModeState.mode],
  );

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',

    // This content is used to create the initial value. It is never referred to again after the first render.
    content: '',
  });
  const { identity: loggedInUserEmail } = useRecoilValue(userData);
  const {
    handleSubmit: handleSendEmail,
    to,
    respondTo,
  } = useRecoilValue(editorEmail);

  const editorRef = useRef<any | null>(null);

  const { onCreatePhoneCallInteractionEvent } =
    useCreatePhoneCallInteractionEvent({ contactId });

  const handleResetEditor = (res: any) => {
    if (!res || !res?.id) return;
    const context = getContext();
    if (context) {
      context.commands.resetContent();
    }
  };
  const handleRespondToEmail = () => {
    const data = prosemirrorNodeToHtml(state.doc);
    if (!handleSendEmail) {
      toast.error('Client error occured while sending the email!');
      return;
    }
    return handleSendEmail(
      data.replace(/(<([^>]+)>)/gi, ''),
      () => console.log(''),
      to,
      respondTo,
    );
  };

  const handleLogPhoneCall = () => {
    const data = prosemirrorNodeToHtml(state.doc);
    const dataToSubmit = {
      appSource: 'Openline',
      sentBy: loggedInUserEmail,
      content: data?.replaceAll(/.src(\S*)/g, '') || '',
      contentType: 'text/html',
    };

    onCreatePhoneCallInteractionEvent(dataToSubmit).then((res) =>
      handleResetEditor(res),
    );
  };

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: 0,
        height: '100%',
        flex: 1,
      }}
      key={editorModeState.mode}
    >
      {editorModeState.mode === EditorMode.Email && <EmailFields />}

      <SocialEditor
        editorRef={editorRef}
        value={''}
        className={'remirror-editor-wrapper-phone-call-editor'}
        manager={manager}
        state={state}
        setState={setState}
        context={getContext()}
        onSubmit={
          editorModeState.mode === EditorMode.Email
            ? handleRespondToEmail
            : handleLogPhoneCall
        }
        submitButtonLabel={
          editorModeState.mode === EditorMode.Email ? 'Send' : 'Log phone call'
        }
      />
    </div>
  );
};
