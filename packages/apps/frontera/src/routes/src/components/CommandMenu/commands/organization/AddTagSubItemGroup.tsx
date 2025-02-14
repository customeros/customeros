import { useMemo, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { EditOrganizationTagUsecase } from '@domain/usecases/command-menu/edit-organization-tag.usecase.ts';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Check } from '@ui/media/icons/Check.tsx';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddTagSubItemGroup = observer(() => {
  const editTagsUsecase = useMemo(() => new EditOrganizationTagUsecase(), []);

  useEffect(() => {
    editTagsUsecase.allowClose();
  }, []);

  return (
    <>
      {editTagsUsecase.tagList?.map((tag) => (
        <CommandSubItem
          key={tag.id}
          icon={<Tag01 />}
          leftLabel='Change tag'
          rightLabel={tag.value.name}
          onSelectAction={() => {
            editTagsUsecase.select(tag.id);
          }}
          rightAccessory={
            editTagsUsecase.organizationTags.has(tag.value.name) ? (
              <Check />
            ) : null
          }
        />
      ))}
    </>
  );
});
