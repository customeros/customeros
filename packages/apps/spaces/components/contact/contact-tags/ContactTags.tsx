import React from 'react';
import { useContactTags } from '../../../hooks/useContact';
import { TagsList, TagListSkeleton } from '../../ui-kit';
export const ContactTags = ({ id }: { id: string }) => {
  const { data, loading, error } = useContactTags({ id });

  if (loading) {
    return <TagListSkeleton />;
  }
  if (error) {
    return null;
  }

  return <TagsList tags={data?.tags ?? []} onTagDelete={() => null} readOnly />;
};
