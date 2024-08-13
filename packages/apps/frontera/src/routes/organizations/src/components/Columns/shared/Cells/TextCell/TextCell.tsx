import React, { ReactNode } from 'react';

interface TextCellProps {
  text: string;
  leftIcon?: ReactNode;
}

export const TextCell = ({ text, leftIcon }: TextCellProps) => {
  if (!text) return <div className='text-gray-400'>Unknown</div>;

  return (
    <div className='overflow-x-hidden overflow-ellipsis flex'>
      {leftIcon && <div className='mr-1'>{leftIcon}</div>}
      {text}
    </div>
  );
};
