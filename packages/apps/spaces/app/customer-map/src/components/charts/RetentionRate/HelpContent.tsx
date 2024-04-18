export const HelpContent = () => {
  return (
    <div className='mt-1'>
      <p className='text-sm font-normal'>
        Retention Rate is the percentage of customers who continue to subscribe
        to your service over a specific period.
      </p>
      <br />
      <p className='text-sm font-normal'>
        To calculate this rate, we look at the number of customers with a
        renewal in the current period and determine what percentage in fact
        renewed their subscription.
      </p>
      <p className='text-sm font-normal'>
        For example, customers with an Annual renewal in January will only be
        included in this metric during the January period as they are not
        eligible for renewal from February to December.
      </p>
      <br />
      <p className='text-sm font-normal'>
        The higher this percentage, the more effectively you are maintaining
        your customer base.
      </p>
    </div>
  );
};
