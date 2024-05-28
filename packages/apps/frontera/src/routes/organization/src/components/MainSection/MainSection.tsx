import { useParams } from 'react-router-dom';

import { UserPresence } from '@shared/components/UserPresence';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const MainSection = ({ children }: { children?: React.ReactNode }) => {
  const organizationId = useParams()?.id as string;

  return (
    <Card
      id='main-section'
      className='flex h-full flex-grow flex-shrink border-none rounded-none flex-col overflow-hidden shadow-none relative bg-gray-25 min-w-[609px] p-0'
    >
      <CardHeader className='px-6 pt-5 pb-2 flex items-center flex-row justify-between'>
        <h1 className='font-semibold text-lg text-gray-700'>Timeline</h1>
        <UserPresence channelName={`organization_presence:${organizationId}`} />
      </CardHeader>
      <CardContent className='p-0 flex flex-col flex-1'>{children}</CardContent>
    </Card>
  );
};
