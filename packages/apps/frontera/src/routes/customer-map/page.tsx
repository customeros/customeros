import { CustomerMap } from './src/components/charts/CustomerMap';
import { ARRBreakdown } from './src/components/charts/ARRBreakdown';
import { NewCustomers } from './src/components/charts/NewCustomers';
import { RevenueAtRisk } from './src/components/charts/RevenueAtRisk';
import { RetentionRate } from './src/components/charts/RetentionRate';
import { TimeToOnboard } from './src/components/charts/TimeToOnboard';
import { MrrPerCustomer } from './src/components/charts/MrrPerCustomer';
import { OnboardingCompletion } from './src/components/charts/OnboardingCompletion';
import { GrossRevenueRetention } from './src/components/charts/GrossRevenueRetention';

export const DashboardPage = () => {
  return (
    <div className='flex flex-col pl-3 pt-4 overflow-y-auto'>
      <div className='flex mb-6'>
        <CustomerMap />
      </div>

      <div className='flex gap-3 mb-3'>
        <MrrPerCustomer />
        <GrossRevenueRetention />
      </div>

      <div className='flex gap-3 mb-3'>
        <ARRBreakdown />
        <RevenueAtRisk />
      </div>

      <div className='flex gap-3 mb-3'>
        <NewCustomers />
        <RetentionRate />
      </div>

      <div className='flex gap-3'>
        <TimeToOnboard />
        <OnboardingCompletion />
      </div>
    </div>
  );
};
