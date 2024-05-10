import { twMerge } from 'tailwind-merge';
import { clsx, type ClassValue } from 'clsx';

export function cn(...args: ClassValue[]) {
  return twMerge(clsx(args));
}
