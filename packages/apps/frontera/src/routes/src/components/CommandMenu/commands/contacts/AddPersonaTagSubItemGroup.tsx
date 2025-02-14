import { useMemo, useEffect } from 'react';

import { EditPersonaTagUsecase } from '@domain/usecases/command-menu/edit-persona-tag.usecase';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Check } from '@ui/media/icons/Check';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddPersonaTagSubItemGroup = () => {
  const usecase = useMemo(() => new EditPersonaTagUsecase(), []);

  useEffect(() => {
    usecase.allowClose();
  }, []);

  return (
    <>
      {usecase.tagList?.map((tag) => (
        <CommandSubItem
          key={tag.id}
          icon={<Tag01 />}
          leftLabel='Add tag'
          rightLabel={tag.value.name}
          onSelectAction={() => {
            usecase.select(tag.id);
          }}
          rightAccessory={
            usecase.contactTags.has(tag.value.name) ? <Check /> : null
          }
        />
      ))}
    </>
  );
};
