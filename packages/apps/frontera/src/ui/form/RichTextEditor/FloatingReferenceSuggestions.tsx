import { FC, useState, useEffect } from 'react';

import {
  useMentionAtom,
  FloatingWrapper,
  MentionAtomNodeAttributes,
} from '@remirror/react';

import { cn } from '@ui/utils/cn';

export const FloatingReferenceSuggestions: FC<{
  tags?: Array<{ id: string; label: string }>;
  mentionOptions?: Array<{ id: string; label: string }>;
}> = ({ tags = [], mentionOptions = [] }) => {
  const [options, setOptions] = useState<MentionAtomNodeAttributes[]>([]);
  const { state, getMenuProps, getItemProps, indexIsHovered, indexIsSelected } =
    useMentionAtom({
      items: options,
      // @ts-expect-error space is not included in types but it's a valid option
      submitKeys: ['Space', 'Enter'],
    });

  useEffect(() => {
    if (!state) return;

    const searchTerm = state.query.full.toLowerCase();
    const options = state.name === 'tag' ? tags : mentionOptions;

    let filteredOptions: { id: string; label: string; hide?: boolean }[] =
      options
        .filter((option) => option.label.toLowerCase().includes(searchTerm))
        .sort()
        .slice(0, 4);

    if (state.name === 'reference' && filteredOptions.length === 0) {
      filteredOptions = [{ id: searchTerm, label: searchTerm, hide: true }];
    }
    setOptions(filteredOptions);
  }, [state]);

  const enabled = Boolean(state);

  return (
    <FloatingWrapper
      placement='auto'
      enabled={enabled}
      positioner='cursor'
      renderOutsideEditor
    >
      <div {...getMenuProps()} className='floating-menu'>
        {enabled &&
          options.map((reference, index) => {
            const isHighlighted = indexIsSelected(index);
            const isHovered = indexIsHovered(index);

            if (reference?.hide) {
              return (
                <div
                  key={`remirror-mention-reference-suggestion-${reference.label}-${reference.id}`}
                  {...getItemProps({
                    item: reference,
                    index,
                  })}
                />
              );
            }

            return (
              <div
                key={`remirror-mention-reference-suggestion-${reference.label}-${reference.id}`}
                className={cn(
                  'floating-menu-option',
                  isHighlighted && 'highlighted',
                  isHovered && 'hovered',
                )}
                {...getItemProps({
                  item: reference,
                  index,
                })}
              >
                {reference.label}
              </div>
            );
          })}
      </div>
    </FloatingWrapper>
  );
};
