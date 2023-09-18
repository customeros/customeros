import React from 'react';
import { Button } from '@ui/form/Button';

interface TagButtonProps {
  onTagSet: () => void;
  tag: string;
}

export const TagButton: React.FC<TagButtonProps> = ({ onTagSet, tag }) => (
  <Button
    size='xs'
    fontSize='md'
    fontWeight='normal'
    color='gray.400'
    variant='unstyled'
    mr={2}
    onClick={onTagSet}
  >
    {`#${tag}`}
  </Button>
);
