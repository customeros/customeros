import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';

import { useRemirror } from '@remirror/react';
import { htmlToProsemirrorNode } from 'remirror';

import { Contact } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getMentionOptionLabel } from '@organization/src/hooks/utils';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { RichEditorBlurHandler } from '@ui/form/RichTextEditor/components/RichEditorBlurHandler';
import { useGetMentionOptionsQuery } from '@organization/src/graphql/getMentionOptions.generated';
import { FloatingReferenceSuggestions } from '@ui/form/RichTextEditor/FloatingReferenceSuggestions';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { logEntryEditorExtensions } from '@organization/src/components/Timeline/TimelineActions/context/extensions';

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
      <Text
        position='relative'
        sx={{
          '& .remirror-theme': { fontSize: 'sm' },
        }}
      >
        <RichTextEditor
          {...remirrorProps}
          placeholder='Log a conversation you had with a customer'
          formId={formId}
          name='content'
          showToolbar={false}
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

            <RichEditorBlurHandler formId={formId} name='content' />
          </>
        </RichTextEditor>
      </Text>
    </>
  );
};
