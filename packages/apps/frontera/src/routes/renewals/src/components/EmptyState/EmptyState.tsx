import { useNavigate } from 'react-router-dom';

import { Button } from '@ui/form/Button/Button';
import { EmptyTable } from '@ui/media/logos/EmptyTable';

import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

export const EmptyState = () => {
  const navigate = useNavigate();

  const options = {
    title: 'No contracts created yet',
    description:
      'Currently, you have not been assigned to any organizations.\n' +
      '\n' +
      'Head to your list of organizations and assign yourself as an owner to one of them.',
    buttonLabel: 'Go to Organizations',
    onClick: () => {
      navigate(`/finder`);
    },
  };

  return (
    <div className='flex items-center justify-center h-full bg-white'>
      <div className='flex flex-col h-[500px] w-[500px]'>
        <div className='flex relative'>
          <EmptyTable className='w-[152px] h-[120px] absolute top-[25%] right-[35%]' />
          <HalfCirclePattern width={500} height={500} />
        </div>
        <div className='flex flex-col text-center items-center top-[5vh] transform translate-y-[-230px]'>
          <p className='text-gray-900 text-base font-semibold'>
            {options.title}
          </p>
          <p className='max-w-[400px] text-sm text-gray-600 my-1'>
            {options.description}
          </p>

          <Button
            variant='outline'
            onClick={options.onClick}
            className='mt-2 min-w-min text-sm'
          >
            {options.buttonLabel}
          </Button>
        </div>
      </div>
    </div>
  );
};
