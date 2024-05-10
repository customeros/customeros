import { AnimatePresence } from 'framer-motion';

export const TabsContainer = ({ children }: { children?: React.ReactNode }) => {
  return (
    <AnimatePresence mode='wait'>
      <div className='flex w-full h-[100%] bg-gray-25 flex-col border-r border-gray-200 overflow-hidden'>
        {children}
      </div>
    </AnimatePresence>
  );
};
