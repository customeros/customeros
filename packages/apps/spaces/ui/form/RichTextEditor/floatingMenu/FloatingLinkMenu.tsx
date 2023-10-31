import { FloatingWrapper, useCurrentSelection } from '@remirror/react';
import { useEffect, useState } from 'react';

import { LinkComponent } from '@ui/form/RichTextEditor/floatingMenu/LinkInput';

export const FloatingLinkToolbar = () => {
  const [isEditing, setIsEditing] = useState(false);
  const { to, from } = useCurrentSelection();

  useEffect(() => {
    if (isEditing && from === to) {
      setIsEditing(false);
    }
    if (!isEditing && from !== to) {
      setIsEditing(true);
    }
  }, [from, to, isEditing]);

  return (
    <>
      <FloatingWrapper
        positioner='selection'
        placement='auto'
        enabled={isEditing}
      >
        {isEditing && <LinkComponent isEditing={isEditing} />}
      </FloatingWrapper>
    </>
  );
};
