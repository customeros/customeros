import { cn } from '@ui/utils/cn';

import { LoadingMessage } from './LoadingMessage.tsx';
import { SimulatedProgress } from './SimulatedProgressBar.tsx';
import logoCustomerOs from '../../../assets/customer-os-small.png';

export const LoadingScreen = ({
  showSplash,
  hide,
  isLoaded,
}: {
  hide: boolean;
  isLoaded: boolean;
  showSplash: boolean;
}) => {
  return (
    <div
      className={cn(
        'absolute flex items-center justify-center top-0 right-0 bottom-0 left-0 z-10 bg-white opacity-0 duration-500 transition-opacity ',
        showSplash && 'opacity-100',
        hide && 'hidden',
      )}
    >
      <div className='w-full flex justify-center items-center flex-col'>
        <div>
          <img src={logoCustomerOs} alt='CustomerOS' width={44} height={44} />
        </div>
        <h1 className='text-md font-medium mt-2'>Please wait...</h1>
        <LoadingMessage />
        <div className='mt-4 w-full max-w-[353px]'>
          <SimulatedProgress accelerate={isLoaded} />
        </div>
      </div>
    </div>
  );
};
