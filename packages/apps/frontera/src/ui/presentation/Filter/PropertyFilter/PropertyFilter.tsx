import { Button } from '@ui/form/Button/Button';

interface PropertyFilterInterface {
  name: string;
}

export const PropertyFilter = ({ name }: PropertyFilterInterface) => {
  return (
    <Button
      size='xs'
      isDisabled
      colorScheme='grayModern'
      className='border-transparent'
    >
      {name || 'Property'}
    </Button>
  );
};
