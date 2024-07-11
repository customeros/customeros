import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { EmptyTable } from '@ui/media/logos/EmptyTable';

import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

export const EmptyState = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const preset = searchParams?.get('preset');

  const currentPreset = store.tableViewDefs
    ?.toArray()
    .find((e) => e.value.id === preset)?.value?.name;
  const handleCreateOrganization = () => {
    store.organizations.create();
  };

  const leadsView = store.tableViewDefs
    ?.toArray()
    .find((e) => e.value.name === 'Leads')?.value.id;

  const options = (() => {
    switch (currentPreset) {
      case 'All orgs':
        return {
          title: "Let's get started",
          description:
            'Start seeing your customer conversations all in one place by adding an organization',
          buttonLabel: 'Add organization',
          dataTest: 'all-orgs-add-org',
          onClick: handleCreateOrganization,
        };
      case 'Contacts':
        return {
          title: 'No contacts created yet',
          description: 'Currently, there are no contacts created yet.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'contacts-go-to-orgs',
          onClick: () => {
            navigate(`/finder`);
          },
        };
      case 'Portfolio':
        return {
          title: 'No organizations assigned to you yet',
          description:
            'Currently, you have not been assigned to any organizations.\n' +
            '\n' +
            'Head to your list of organizations and assign yourself as an owner to one of them.',
          buttonLabel: 'Go to All orgs',
          dataTest: 'portfolio-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder`);
          },
        };
      case 'Customers':
        return {
          title: 'Who will be first?',
          description:
            'No customers here yet. You can change prospects into customers by changing their relationship status in the About section.',
          buttonLabel: 'Go to All orgs',
          dataTest: 'customers-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder`);
          },
        };
      case 'Leads':
        return {
          title: 'Lead-free zone',
          description:
            'We’re on the lookout for new leads. Once we find them, they will appear here, or automatically qualified by your Ideal Company Profile.',
        };
      case 'Targets':
        return {
          title: 'Bullseye pending',
          description:
            'We’re sorting through your leads using your Ideal Company Profile . Once qualified, they will automatically show up here as Targets.',
          buttonLabel: 'Go to Leads',
          dataTest: 'targets-go-to-leads',
          onClick: () => {
            navigate(`/finder?preset=${leadsView}`);
          },
        };
      case 'Churn':
        return {
          title: 'Smooth sailing',
          description:
            'Seems like your customers are loyal! No one has churned yet. Keep up the strong relationships.',
          buttonLabel: 'Go to All orgs',
          dataTest: 'churn-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder`);
          },
        };
      default:
        return {
          title: "We couldn't find any organizations",
          description:
            'Organizations whose relationship is set to Customer, will appear here. You can change this in an organization’s About section.',
          buttonLabel: 'Go to All orgs',
          dataTest: 'go-to-all-orgs',
          onClick: () => {
            navigate(`/finder`);
          },
        };
    }
  })();

  return (
    <div className='flex items-center justify-center h-full bg-white'>
      <div className='flex flex-col h-[500px] w-[500px]'>
        <div className='flex relative'>
          <EmptyTable className='w-[152px] h-[120px] absolute top-[25%] right-[35%]' />
          <HalfCirclePattern height={500} width={500} />
        </div>
        <div className='flex flex-col text-center items-center top-[5vh] transform translate-y-[-230px]'>
          <p className='text-gray-900 text-md font-semibold'>{options.title}</p>
          <p className='max-w-[400px] text-sm text-gray-600 my-1'>
            {options.description}
          </p>

          {currentPreset !== 'Leads' && (
            <Button
              onClick={options.onClick}
              className='mt-4 min-w-min text-sm'
              data-test={options.dataTest}
              variant='outline'
            >
              {options.buttonLabel}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
});
