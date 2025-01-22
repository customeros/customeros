import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';

export const Header = () => {
  return (
    <div className='w-full border-b border-b-gray-200 px-3 py-2 flex justify-between items-center'>
      <p className='text-md font-medium'>Agents</p>

      <Button size='xs' leftIcon={<Plus />} colorScheme='primary'>
        New agent
      </Button>
    </div>
  );
};
