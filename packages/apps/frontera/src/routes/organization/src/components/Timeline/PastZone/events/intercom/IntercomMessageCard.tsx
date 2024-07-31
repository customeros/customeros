import React, { PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import { User01 } from '@ui/media/icons/User01';
import { Avatar } from '@ui/media/Avatar/Avatar';
import { Intercom } from '@ui/media/icons/Intercom';
import { ViewInExternalAppButton } from '@ui/form/Button';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';

interface IntercomMessageCardProps extends PropsWithChildren {
  name: string;
  date: string;
  content: string;
  className?: string;
  onClick?: () => void;
  sourceUrl?: string | null;
  showDateOnHover?: boolean;
  profilePhotoUrl?: null | string;
}

export const IntercomMessageCard: React.FC<IntercomMessageCardProps> = ({
  name,
  sourceUrl,
  profilePhotoUrl,
  content,
  onClick,
  children,
  className,
  date,
  showDateOnHover,
}) => {
  return (
    <>
      <Card
        onClick={() => onClick?.()}
        className={cn(
          className,
          onClick ? 'cursor-pointer' : '',
          'text-sm flex shadow-xs border border-gray-200 bg-white [intercom-stub-date]:hover:text-gray-500 max-w-[549px]',
        )}
      >
        <CardContent className='p-3 overflow-hidden w-full'>
          <div className='flex gap-3 flex-1'>
            <Avatar
              size='md'
              name={name}
              variant='roundedSquare'
              src={profilePhotoUrl || undefined}
              icon={<User01 height='1.8rem' color='gray.500' />}
              className={cn(
                profilePhotoUrl ? 'border-none' : 'border border-primary-200',
              )}
            />
            <div
              className={cn(
                'flex flex-1 flex-col relative',
                showDateOnHover ? 'max-w-[470px]' : 'max-w-[408px]',
              )}
            >
              <div className='flex justify-between flex-1'>
                <div className='flex items-baseline'>
                  <p className='text-gray-700 font-semibold'>{name}</p>
                  <p
                    className={cn(
                      showDateOnHover ? 'transparent' : 'text-gray-500',
                      'ml-2 text-xs intercom-stub-date',
                    )}
                  >
                    {date}
                  </p>
                </div>

                <ViewInExternalAppButton
                  url={sourceUrl}
                  icon={
                    <div className='flex items-center justify-center'>
                      <Intercom className='h-10' />
                    </div>
                  }
                />
              </div>
              <HtmlContentRenderer
                htmlContent={content}
                noOfLines={showDateOnHover ? 4 : undefined}
                pointerEvents={showDateOnHover ? 'none' : 'initial'}
              />
              {children}
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
