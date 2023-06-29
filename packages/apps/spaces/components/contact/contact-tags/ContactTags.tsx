import React from 'react';
import { useRemoveTagFromContact } from '@spaces/hooks/useContact';
import { TagsList } from '@spaces/atoms/tags';
import { ContactTagsEdit } from './ContactTagsEdit';
import { TagFragment } from '@spaces/graphql';
export const ContactTags = ({
  id,
  mode,
  tags,
}: {
  id: string;
  mode: 'PREVIEW' | 'EDIT';
  tags?: Array<TagFragment> | null;
}) => {
  const { onRemoveTagFromContact } = useRemoveTagFromContact({ contactId: id });

  return (
    <section style={{ display: 'flex' }}>
      <TagsList
        tags={tags ?? []}
        onTagDelete={(id) => onRemoveTagFromContact({ tagId: id })}
        readOnly={mode === 'PREVIEW'}
      >
        {mode === 'EDIT' && (
          <ContactTagsEdit contactId={id} contactTags={tags || []} />
        )}
      </TagsList>
    </section>
  );
};
