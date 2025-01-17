import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

interface DomainCellProps {
  organizationId: string;
}

export const DomainsCell = observer(({ organizationId }: DomainCellProps) => {
  const store = useStore();
  const organization = store.organizations.getById(organizationId);

  const domains = organization?.primaryDomains;

  return (
    <div className='flex items-center cursor-pointer'>
      <p className='text-gray-700  truncate'>
        {domains?.length ? (
          <span className='truncate'>
            {domains.map((domain, index) => (
              <a
                target='_blank'
                rel='noopener noreferrer'
                href={`https://${domain}`}
                className='hover:underline'
                key={`primary-domain-${domain}`}
              >
                {domain}
                {index < domains.length - 1 && ', '}
              </a>
            ))}
          </span>
        ) : organization?.isEnriching ? (
          <span
            className='text-gray-400'
            data-test='organization-website-in-all-orgs-table'
          >
            Enriching...
          </span>
        ) : (
          <span
            className='text-gray-400'
            data-test='organization-website-in-all-orgs-table'
          >
            Not set
          </span>
        )}
      </p>
    </div>
  );
});
