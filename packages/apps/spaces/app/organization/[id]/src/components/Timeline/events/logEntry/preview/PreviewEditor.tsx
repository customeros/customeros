import React, { useEffect } from 'react';
import { Text } from '@ui/typography/Text';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { FloatingReferenceSuggestions } from '@ui/form/RichTextEditor/FloatingReferenceSuggestions';
import { useRemirror } from '@remirror/react';
import { logEntryEditorExtensions } from '@organization/src/components/Timeline/TimelineActions/context/extensions';
import { RichEditorBlurHandler } from '@ui/form/RichTextEditor/components/RichEditorBlurHandler';
import { htmlToProsemirrorNode } from 'remirror';
import { useGetMentionOptionsQuery } from '@organization/src/graphql/getMentionOptions.generated';
import { useParams } from 'next/navigation';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { getMentionOptionLabel } from '@organization/src/components/Timeline/events/utils';
import { Contact } from '@graphql/types';

export const PreviewEditor: React.FC<{
  formId: string;
  initialContent: string;
  tags?: Array<{ label: string; value: string }>;
}> = ({ formId, initialContent, tags }) => {
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
    .filter((e) => Boolean(e.label)) as { label: string; id: string }[];

  return (
    <>
      <Text sx={{ '& .remirror-theme': { fontSize: 'sm' } }}>
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
            <RichEditorBlurHandler formId={formId} name='content' />
          </>
        </RichTextEditor>
      </Text>
    </>
  );
};
