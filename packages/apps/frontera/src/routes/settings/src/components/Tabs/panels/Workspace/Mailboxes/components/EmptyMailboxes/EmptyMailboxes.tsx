import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { Inbox01 } from '@ui/media/icons/Inbox01';

interface EmptyMailboxesProps {
  onUpdate: () => void;
}

export const EmptyMailboxes = ({ onUpdate }: EmptyMailboxesProps) => {
  const handleButtonClick = () => {
    onUpdate();
  };

  return (
    <div className=' border-r-[1px] max-w-[600px] h-full'>
      <div
        className='flex flex-col h-[400px] bg-no-repeat bg-cover max-w-[600px]'
        style={{
          backgroundImage: `url(/backgrounds/organization/dotted-bg-pattern.svg)`,
          backgroundPositionX: 'center',
          backgroundPositionY: 'center',
        }}
      >
        <div className='flex flex-col items-center '>
          <FeaturedIcon colorScheme='gray' className='mt-[100px]'>
            <Inbox01 />
          </FeaturedIcon>
          <div className='px-[128px]'>
            <p className='mt-6 font-medium text-center'>
              Set up outbound mailboxes
            </p>
            <p className='[300px] text-center text-sm'>
              Start sending emails in Flows by setting up outbound mailboxes for
              your workspace
            </p>
          </div>

          <Button
            className='mt-6'
            colorScheme='primary'
            onClick={handleButtonClick}
          >
            Set up mailboxes
          </Button>
        </div>
      </div>
    </div>
  );
};
