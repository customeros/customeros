import React, { useRef, useMemo, useState, forwardRef } from 'react';

import { useMergeRefs } from 'rooks';

import { cn } from '@ui/utils/cn';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';

interface ReminderPostitProps extends React.HTMLAttributes<HTMLDivElement> {
  owner?: string;
  className?: string;
  isFocused?: boolean;
  isMutating?: boolean;
  onClickOutside?: (e: Event) => void;
}

const rotations = ['rotate(2deg)', 'rotate(-2deg)', 'rotate(0deg)'];
const rgadients = [
  'linear-gradient(to top left, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.06) 100%)',
  'linear-gradient(to top right, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.06) 100%)',
  'linear-gradient(to bottom left, rgba(196, 196, 196, 0.00) 20%, rgba(0, 0, 0, 0.03) 100%)',
];

const getRandomStyles = () => {
  const index = Math.floor(Math.random() * rotations.length);

  return [rotations[index], rgadients[index]];
};

export const ReminderPostit = forwardRef<HTMLDivElement, ReminderPostitProps>(
  (
    {
      owner,
      children,
      isFocused,
      className,
      isMutating,
      onClickOutside = () => undefined,
      ...rest
    },
    ref,
  ) => {
    const _ref = useRef(null);
    const [isHovered, setIsHovered] = useState(false);
    const [rotation, gradient] = useMemo(() => getRandomStyles(), []);

    const mergedRef = useMergeRefs(_ref, ref);

    useOutsideClick({ ref: _ref, handler: onClickOutside });

    return (
      <div
        ref={mergedRef}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        className={cn(
          className,
          isMutating
            ? 'pointer-events-none animate-pulseOpacity duration-75 ease-in-out'
            : 'pointer-events-auto',
          'flex relative w-[321px] m-6 mt-2',
        )}
        {...rest}
      >
        <div
          style={{ transform: isFocused || isHovered ? 'unset' : rotation }}
          className={cn(
            isFocused || isHovered ? 'blur-[7px]' : 'blur-[3px]',
            isFocused || isHovered
              ? 'bg-[rgba(0,0,0,0.2)]'
              : 'bg-[rgba(0,0,0,0.07)]',
            'flex w-[calc(100%-10px)] h-[calc(100%-28px)] bottom-[-4px] left-[5px] absolute transition-all duration-100 ease-in-out',
          )}
        />
        <div className='flex w-full z-[1] bg-[#FEFCBF] flex-col'>
          <div
            className='h-6 w-full items-center'
            style={{
              backgroundImage: `${gradient}`,
            }}
          >
            {owner && (
              <span className='text-xs text-gray-500 pl-4 font-normal pt-3'>
                {owner} added
              </span>
            )}
          </div>
          {children}
        </div>
      </div>
    );
  },
);
