import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { Sale03 } from '@ui/media/icons/Sale03';
import { Inbox01 } from '@ui/media/icons/Inbox01';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

interface EmptyMailboxesProps {
  onUpdate: () => void;
}

export const EmptyMailboxes = ({ onUpdate }: EmptyMailboxesProps) => {
  const handleButtonClick = () => {
    onUpdate();
  };

  return (
    <div className=' border-r-[1px] max-w-[418px] h-full'>
      <div className='flex flex-col pl-6 pr-5 pt-10'>
        <FeaturedIcon colorScheme='gray'>
          <Inbox01 />
        </FeaturedIcon>
        <div className='mt-4 gap-1'>
          <p className='font-medium'>Outbound mailboxes in 3 minutes</p>
          <p className='text-sm mt-2'>
            Set up your outbound email in minutes. Fully automated, secure, and
            built to scale—send 10,000+ emails per month with professional-grade
            tools that do the hard work for you.
          </p>
          <p className='text-sm font-medium pt-3 mb-1'>
            The Starter bundle gives you:
          </p>
          <div className='px-3 gap-1'>
            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>20 mailboxes</span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>Automatic email setup</span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>
                Built-in DKIM, DMARC, and SPF protection
              </span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>Email and domain forwarding</span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>Custom domain tracking</span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>
                Continuous mailbox reputation monitoring
              </span>
            </div>

            <div className='flex items-center gap-2 py-1'>
              <CheckCircle className='text-success-500' />
              <span className='text-sm'>
                Automated sequences with CustomerOS Flows
              </span>
            </div>
          </div>
        </div>

        <div className='flex items-center py-2 px-3 border rounded-md border-grayModern-200 mt-4'>
          <Sale03 className='text-gray-500 mr-2' />
          <span className='text-sm'>
            Get the <span className='font-medium'>Starter bundle</span> for{' '}
            <span className='text-success-500 font-medium'>$199.99/month</span>
          </span>
        </div>

        <Button
          className='mt-5 mx-3'
          colorScheme='primary'
          dataTest='set-up-mailboxes'
          onClick={handleButtonClick}
          rightIcon={<ChevronRight />}
        >
          Set up your mailboxes
        </Button>
      </div>
    </div>
  );
};
