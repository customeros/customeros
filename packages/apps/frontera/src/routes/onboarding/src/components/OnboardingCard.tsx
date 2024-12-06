import { ReactElement } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { Play } from '@ui/media/icons/Play.tsx';
import { Clock } from '@ui/media/icons/Clock.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { CheckVerified03 } from '@ui/media/icons/CheckVerified03.tsx';
import {
  Card,
  CardFooter,
  CardHeader,
  CardContent,
} from '@ui/presentation/Card/Card.tsx';

export const OnboardingCard = ({
  title,
  titleIcon,
  contentTitle,
  contentIcon,
  description,
  isFeatured,
  timeInMinutes,
  subtitle,
  isCompleted,
  onClick,
}: {
  title: string;
  subtitle: string;
  isFeatured: boolean;
  description: string;
  contentTitle: string;
  onClick?: () => void;
  isCompleted?: boolean;
  timeInMinutes: string;
  titleIcon: ReactElement;
  contentIcon: ReactElement;
}) => {
  return (
    <Card
      className={cn('p-4 max-w-[332px] shadow-md relative bg-white', {
        'bg-grayModern-50': isCompleted,
        'p-8 max-w-[370px]': isFeatured,
      })}
    >
      <CardHeader className='border-b border-grayModern-200 mb-3 pb-3'>
        <h1
          className={cn(
            'text-base font-medium inline-flex items-center gap-3',
            {
              'text-xl font-bold': isFeatured,
              'text-primary-700': !isCompleted,
            },
          )}
        >
          {titleIcon}
          {title}
        </h1>
        <p className={cn('ml-7 text-sm', { 'ml-9 text-base': isFeatured })}>
          {subtitle}
        </p>
      </CardHeader>
      <CardContent className='p-0 mb-3'>
        <h2
          className={cn('font-medium text-sm', {
            'text-base': isFeatured,
          })}
        >
          <span className={cn('mr-3 text-gray-500', { 'mr-4': isFeatured })}>
            {contentIcon}
          </span>
          {contentTitle}
        </h2>
        <p className={cn('ml-7 text-sm', { 'ml-9 text-base': isFeatured })}>
          {description}
        </p>
      </CardContent>
      <CardFooter className='p-0 mt-4'>
        {isCompleted && (
          <div
            className={cn(
              'w-full bg-grayModern-100 p-1.5 rounded flex items-center justify-center gap-2 text-sm font-medium',
              {
                'h-[36px]': isFeatured,
              },
            )}
          >
            <CheckVerified03 />
            Done
          </div>
        )}
        {!isCompleted && (
          <div className='flex justify-between w-full'>
            <div className='flex items-center gap-2'>
              <Clock className='text-gray-500' />
              <span className='text-gray-500 text-sm'>
                Time: {timeInMinutes} minutes
              </span>
            </div>
            <Button
              onClick={onClick}
              leftIcon={<Play />}
              size={isFeatured ? 'sm' : 'xs'}
              colorScheme={isFeatured ? 'primary' : 'gray'}
            >
              Get Started
            </Button>
          </div>
        )}
      </CardFooter>
    </Card>
  );
};
