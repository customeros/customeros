import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { useRemirror } from '@remirror/react';
import { htmlToProsemirrorNode } from 'remirror';

import { Contact } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getMentionOptionLabel } from '@organization/hooks/utils';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { useGetMentionOptionsQuery } from '@organization/graphql/getMentionOptions.generated';
import { RichEditorBlurHandler } from '@ui/form/RichTextEditor/components/RichEditorBlurHandler';
import { FloatingReferenceSuggestions } from '@ui/form/RichTextEditor/FloatingReferenceSuggestions';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { logEntryEditorExtensions } from '@organization/components/Timeline/FutureZone/TimelineActions/context/extensions';

export const PreviewEditor: React.FC<{
  formId: string;
  onClose: () => void;
  initialContent: string;
  tags?: Array<{ label: string; value: string }>;
}> = ({ formId, initialContent, tags, onClose }) => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const { data } = useGetMentionOptionsQuery(client, {
    id,
  });
  const remirrorProps = useRemirror({
    extensions: logEntryEditorExtensions,
  });

  useEffect(() => {
    const prosemirrorNodeValue = htmlToProsemirrorNode({
      schema: remirrorProps.state.schema,
      content: `${initialContent}`,
    });

    remirrorProps.getContext()?.setContent(prosemirrorNodeValue);
  }, [initialContent]);

  const mentionOptions = (data?.organization?.contacts?.content ?? [])
    .map((e) => ({ label: getMentionOptionLabel(e as Contact), id: e.id }))
    .filter((e) => Boolean(e.label)) as { id: string; label: string }[];

  return (
    <>
      <p className='text-sm relative [remirror-theme]:text-sm'>
        <RichTextEditor
          {...remirrorProps}
          name='content'
          formId={formId}
          showToolbar={false}
          placeholder='Log a conversation you had with a customer'
        >
          <>
            <FloatingReferenceSuggestions
              mentionOptions={mentionOptions}
              tags={tags?.map((e: { label: string; value: string }) => ({
                label: e.label,
                id: e.value,
              }))}
            />
            <KeymapperClose onClose={onClose} />

            <RichEditorBlurHandler name='content' formId={formId} />
          </>
        </RichTextEditor>
      </p>
    </>
  );
};
