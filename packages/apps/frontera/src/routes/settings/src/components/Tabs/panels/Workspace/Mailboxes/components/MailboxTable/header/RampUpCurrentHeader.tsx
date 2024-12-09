import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';

export const RampUpCurrentHeader = () => {
  const { open, onOpen, onClose } = useDisclosure();

  return (
    <>
      <div className='flex items-center gap-1'>
        <span className='text-sm'>Daily Email Limit </span>
        <InfoCircle
          onClick={onOpen}
          className='size-3 text-gray-500 cursor-pointer hover:text-gray-700'
        />
      </div>

      <InfoDialog
        isOpen={open}
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
        label='Email best practices'
        body={
          <div className='space-y-4'>
            <p className='text-sm'>
              To maintain deliverability and avoid spam filters, we auto-warm
              and rotate mailboxes, limiting each to a maximum of 40 emails per
              day.
            </p>
            <p className='text-sm'>
              Using two mailboxes per domain has proven especially effective.
            </p>
          </div>
        }
      />
    </>
  );
};
