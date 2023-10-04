import { FC, useEffect, useState } from 'react';
import { cx } from '@remirror/core';
import {
  FloatingWrapper,
  MentionAtomNodeAttributes,
  useMentionAtom,
} from '@remirror/react';
import { Box } from '@ui/layout/Box';

export const FloatingReferenceSuggestions: FC<{
  tags?: Array<{ label: string; id: string }>;
  mentionOptions?: Array<{ label: string; id: string }>;
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

    let filteredOptions: { label: string; id: string; hide?: boolean }[] =
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
      positioner='cursor'
      enabled={enabled}
      placement='auto'
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
              <Box
                key={`remirror-mention-reference-suggestion-${reference.label}-${reference.id}`}
                className={cx(
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
              </Box>
            );
          })}
      </div>
    </FloatingWrapper>
  );
};
