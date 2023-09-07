import { useEffect, useState } from 'react';
import { cx } from '@remirror/core';
import {
  FloatingWrapper,
  MentionAtomNodeAttributes,
  useMentionAtom,
} from '@remirror/react';

const commonTags = ['meeting', 'call', 'voicemail', 'email', 'text-message'];
const tags = commonTags.map((label) => ({ label, id: label }));

export const TagSuggestor: React.FC = () => {
  const [options, setOptions] = useState<MentionAtomNodeAttributes[]>([]);
  const { state, getMenuProps, getItemProps, indexIsHovered, indexIsSelected } =
    useMentionAtom({
      items: options,
      // @ts-expect-error explain
      submitKeys: ['Space', 'Enter'],
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
          options.map((user, index) => {
            const isHighlighted = indexIsSelected(index);
            const isHovered = indexIsHovered(index);

            return (
              <div
                key={user.id}
                className={cx(
                  'floating-menu-option',
                  isHighlighted && 'highlighted',
                  isHovered && 'hovered',
                )}
                {...getItemProps({
                  item: user,
                  index,
                })}
              >
                {user.label}
              </div>
            );
          })}
      </div>
    </FloatingWrapper>
  );
};
