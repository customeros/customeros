import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';

export const Header = () => {
  return (
    <div className='w-full border-b border0b-gray-500 px-3 py-2 flex justify-between items-center'>
      <p className='text-md font-medium'>Agents</p>

      <Button size='xs' leftIcon={<Plus />} colorScheme='primary'>
        New agent
      </Button>
    </div>
  );
};
