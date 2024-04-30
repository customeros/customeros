'use client';
import React from 'react';

export const TabsContainer = ({ children }: { children?: React.ReactNode }) => {
  return (
    <div className='flex min-w-[400px] h-full bg-gray-25 flex-col border-r border-gray-200 overflow-hidden'>
      {children}
    </div>
  );
};
