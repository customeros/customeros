export const HelpContent = () => {
  return (
    <div className='mt-1'>
      <p className='text-sm font-normal'>
        Gross Revenue Retention (GRR) tells you what percentage of revenue you
        keep from all your customers over their lifetime.
      </p>
      <br />
      <p className='text-sm font-normal'>
        To determine this, we compare the Monthly Recurring Revenue (MRR) from
        the current period (minus any up-sells and cross-sells) with the initial
        contracted MRR.
      </p>
      <br />
      <p className='text-sm font-normal'>
        This comparison is then expressed as a percentage, with a maximum value
        of 100% indicating that all original revenue has been retained.
      </p>
    </div>
  );
};
