import React, { useState } from 'react';
import {
  useContactTags,
  useRemoveTagFromContact,
} from '../../../hooks/useContact';
import { TagsList, TagListSkeleton } from '../../ui-kit';
import { ContactTagsEdit } from './ContactTagsEdit';
export const ContactTags = ({
  id,
  mode,
}: {
  id: string;
  mode: 'PREVIEW' | 'EDIT';
}) => {
  const { data, loading, error } = useContactTags({ id });
  const { onRemoveTagFromContact } = useRemoveTagFromContact({ contactId: id });

  if (loading) {
    return <TagListSkeleton />;
  }
  if (error) {
    return null;
  }

  return (
    <section style={{ display: 'flex' }}>
      <TagsList
        tags={data?.tags ?? []}
        onTagDelete={(id) => onRemoveTagFromContact({ tagId: id })}
        readOnly={mode === 'PREVIEW'}
      >
        {mode === 'EDIT' && (
          <ContactTagsEdit contactId={id} contactTags={data?.tags || []} />
        )}
      </TagsList>
    </section>
  );
};
