import React from 'react';

interface TextCellProps {
  text: string;
}

export const TextCell = ({ text }: TextCellProps) => {
  if (!text) return <div className='text-gray-400'>Unknown</div>;

  return <div className='overflow-x-hidden overflow-ellipsis'>{text}</div>;
};
