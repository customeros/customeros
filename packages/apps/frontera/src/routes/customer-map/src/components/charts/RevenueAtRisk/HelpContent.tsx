export const HelpContent = () => {
  return (
    <div className='mt-1'>
      <p className='text-base font-normal'>
        Revenue at risk shows the forecasted revenue from customers whose
        renewal likelihood is rated medium, low or zero in the current period.
      </p>
      <br />
      <p className='text-base font-normal'>
        In contrast, the high confidence segment shows the forecasted revenue
        with a high likelihood to renew.
      </p>
    </div>
  );
};
