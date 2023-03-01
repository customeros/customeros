import React from 'react';
import { TagInput } from '../../ui-kit';
import { useCreateTag, useDeleteTag, useTags } from '../../../hooks/useTags';
import {
  useAddTagToContact,
  useRemoveTagFromContact,
} from '../../../hooks/useContact';

interface ContactTagsEditProps {
  contactId: string;
  contactTags: Array<any>;
}

export const ContactTagsEdit: React.FC<ContactTagsEditProps> = ({
  contactId,
  contactTags,
}) => {
  const {
    tags: tagOptions,
    loading: tagsLoading,
    error: tagsError,
  } = useTags();

  const { onRemoveTagFromContact } = useRemoveTagFromContact({ contactId });
  const { onAddTagToContact } = useAddTagToContact({ contactId });
  const { onCreateTag } = useCreateTag();
  const { onDeleteTag } = useDeleteTag();

  const handleCreateTagForContact = async (name: string) => {
    try {
      const newTag = await onCreateTag({ name });
      if (newTag?.id) {
        return onAddTagToContact({ tagId: newTag?.id, contactId });
      }
    } catch (e) {
      console.log(e);
    }
  };

  if (tagsLoading) {
    return null;
  }
  if (tagsError) {
    return <>ERROR</>;
  }

  return (
    <TagInput
      onNewTag={handleCreateTagForContact}
      onTagChange={() => null}
      onTagRemove={(tagId) => onRemoveTagFromContact({ tagId })}
      tags={contactTags ?? []}
      options={tagsLoading || !tagOptions ? [] : tagOptions}
      onSetTags={() => null}
      onTagSelect={({ id }) => onAddTagToContact({ tagId: id, contactId })}
      onTagDelete={onDeleteTag}
    />
  );
};
