export const MailboxStatus = ({
  hasSenders,
  totalMailboxes,
}: {
  hasSenders: boolean;
  totalMailboxes: number;
}) => {
  const renderContent = () => {
    if (!hasSenders) {
      return 'Add one or more senders to start sending emails in this flow';
    }

    if (totalMailboxes > 0) {
      return (
        <>
          You have{' '}
          <span className='font-medium'>
            {totalMailboxes} {totalMailboxes === 1 ? 'mailbox' : 'mailboxes'}
          </span>{' '}
          available allowing you to send up to{' '}
          <span className='font-medium'>
            {totalMailboxes * 6} emails per day.
          </span>
        </>
      );
    }

    return "You haven't set up any mailboxes yet. Add some to start sending emails.";
  };

  return <p className='text-sm'>{renderContent()}</p>;
};
