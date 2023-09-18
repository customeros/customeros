import { FC, useEffect, useState } from 'react';
import { cx } from '@remirror/core';
import {
  FloatingWrapper,
  MentionAtomNodeAttributes,
  useMentionAtom,
} from '@remirror/react';

export const TagSuggestor: FC<{
  tags: Array<{ label: string; id: string }>;
}> = ({ tags = [] }) => {
  const [options, setOptions] = useState<MentionAtomNodeAttributes[]>(tags);

  const { state, getMenuProps, getItemProps, indexIsHovered, indexIsSelected } =
    useMentionAtom({
      items: options,
      // @ts-expect-error space is not included in types but it's a valid option
      submitKeys: ['Space'],
    });

  useEffect(() => {
    if (!state) {
      return;
    }

    const searchTerm = state.query.full.toLowerCase();

    const filteredOptions = tags
      .filter((tag) => tag.label.toLowerCase().includes(searchTerm))
      .sort()
      .slice(0, 5);

    if (filteredOptions.length > 0) {
      setOptions(filteredOptions);
    }
    if (filteredOptions.length === 0) {
      setOptions([{ id: searchTerm, label: searchTerm }]);
    }
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
          options.map((tag, index) => {
            const isHighlighted = indexIsSelected(index);
            const isHovered = indexIsHovered(index);

            return (
              <div
                key={`remirror-mention-tag-suggestion-${tag.label}-${tag.id}`}
                className={cx(
                  'floating-menu-option',
                  isHighlighted && 'highlighted',
                  isHovered && 'hovered',
                )}
                {...getItemProps({
                  item: tag,
                  index,
                })}
              >
                {tag.label}
              </div>
            );
          })}
      </div>
    </FloatingWrapper>
  );
};
