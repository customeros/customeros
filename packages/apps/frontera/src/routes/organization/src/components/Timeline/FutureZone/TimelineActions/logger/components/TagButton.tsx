import { Button } from '@ui/form/Button/Button';

interface TagButtonProps {
  tag: string;
  onTagSet: () => void;
}

export const TagButton = ({ onTagSet, tag }: TagButtonProps) => (
  <Button
    size='xs'
    color='gray.400'
    onClick={onTagSet}
    className='text-gray-400 mr-2 leading-4'
  >
    {`#${tag}`}
  </Button>
);
