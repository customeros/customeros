'use client';

import { cn } from '@ui/utils/cn';

interface DotProps {
  colorScheme?: string;
}

export const Dot = ({ colorScheme, ...props }: DotProps) => {
  const colors = (colorScheme: string) => {
    switch (colorScheme) {
      case 'gray':
        return 'bg-gray-500';
      case 'error':
        return 'bg-error-500';
      case 'success':
        return 'bg-success-500';
      case 'warning':
        return 'bg-warning-500';
      default:
        return 'text-gray-500';
    }
  };

  return (
    <div
      className={cn(colors(colorScheme || 'gray'), 'size-[10px] rounded-full')}
      {...props}
    />
  );
};
