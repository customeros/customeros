import { Link } from 'react-router-dom';

import { ExternalSystem } from '@graphql/types';
import { Hubspot } from '@ui/media/logos/Hubspot';
import { getExternalUrl } from '@utils/getExternalLink';
import { Salesforce } from '@ui/media/logos/Salesforce';

interface LogEntryExternalLinkProps {
  externalLink: ExternalSystem;
}

const getIcon = (type: string) => {
  switch (type) {
    case 'SALESFORCE':
      return <Salesforce className='size-5 mr-2' />;
    case 'HUBSPOT':
      return <Hubspot className='size-5 mr-2' />;

    default:
      return '';
  }
};

export const LogEntryExternalLink = ({
  externalLink,
}: LogEntryExternalLinkProps) => {
  const icon = (() => getIcon(externalLink.type))();
  const link = getExternalUrl(`${externalLink?.externalUrl}`);

  return (
    <div className='mt-4'>
      <Link to={link}>
        {icon}
        View in {externalLink.type.toLowerCase()}
      </Link>
    </div>
  );
};
