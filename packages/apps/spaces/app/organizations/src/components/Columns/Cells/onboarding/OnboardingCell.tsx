import { OnboardingStatus } from '@graphql/types';

interface OnboardingCellProps {
  status?: OnboardingStatus;
}

export const OnboardingCell = ({
  status = OnboardingStatus.NotApplicable,
}: OnboardingCellProps) => {
  return (
    <>
      <div>{status}</div>
    </>
  );
};
