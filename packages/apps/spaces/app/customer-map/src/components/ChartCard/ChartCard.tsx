'use client';
import { ReactNode, PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import { useDisclosure } from '@ui/utils';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog/InfoDialog2';

import { HelpButton } from '../HelpButton';

interface ChartCardProps extends React.HtmlHTMLAttributes<HTMLDivElement> {
  stat?: string;
  title: string;
  hasData?: boolean;
  renderSubStat?: () => ReactNode;
  renderHelpContent?: () => ReactNode;
}

export const ChartCard = ({
  stat,
  title,
  hasData,
  children,
  renderSubStat,
  renderHelpContent,
  ...props
}: PropsWithChildren<ChartCardProps>) => {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <Card
        className='w-full rounded-lg shadow-none border border-gray-200 group'
        {...props}
      >
        <CardHeader className='pb-0 pt-4 px-6'>
          <div className='flex gap-2 items-center'>
            <p className='text-lg font-normal'>{title}</p>
            {true && <HelpButton isOpen={isOpen} onOpen={onOpen} />}
          </div>
          {stat && (
            <h2
              className={cn(
                hasData ? 'text-gray-700' : 'text-lg text-gray-400',
                'text-3xl font-semibold',
              )}
            >
              {hasData ? stat : 'No data yet'}
            </h2>
          )}
          {hasData && renderSubStat && renderSubStat?.()}
        </CardHeader>
        <CardContent className='px-6 pb-6'>{children}</CardContent>
      </Card>

      <InfoDialog
        label={title}
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
      >
        {renderHelpContent?.()}
      </InfoDialog>
    </>
  );
};
