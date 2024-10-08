import { isValidElement } from 'react';

import { Button } from '@ui/form/Button/Button';

interface PropertyFilterInterface {
  name: string;
  icon: React.ReactNode;
}

export const PropertyFilter = ({ name, icon }: PropertyFilterInterface) => {
  return (
    <Button
      size='xs'
      variant='outline'
      colorScheme='grayModern'
      leftIcon={isValidElement(icon) ? icon : undefined}
      className='cursor-not-allowed font-normal bg-white focus:bg-white hover:bg-white'
    >
      {name || 'Property'}
    </Button>
  );
};
