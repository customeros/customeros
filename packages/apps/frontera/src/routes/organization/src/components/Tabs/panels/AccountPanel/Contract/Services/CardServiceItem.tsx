import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

function getBilledTypeLabel(billedType: BilledType): string {
  switch (billedType) {
    case BilledType.Annually:
      return '/year';
    case BilledType.Monthly:
      return '/month';
    case BilledType.Quarterly:
      return '/quarter';
    default:
      return '';
  }
}

export const CardServiceItem = ({
  data,
  onOpen,
  currency,
}: {
  data: ServiceLineItem;
  currency?: string | null;
  onOpen: (props: ServiceLineItem) => void;
}) => {
  return (
    <>
      <div
        className='flex w-full justify-between cursor-pointer text-sm focus:outline-none'
        onClick={() => onOpen(data)}
      >
        {data.description && <p>{data.description}</p>}
        <div className='flex justify-between'>
          <p>
            {data.quantity}
            <span className='text-sm mx-1'>×</span>

            {formatCurrency(data.price ?? 0, 2, currency || 'USD')}
            {getBilledTypeLabel(data.billingCycle)}
          </p>
        </div>
      </div>
    </>
  );
};
