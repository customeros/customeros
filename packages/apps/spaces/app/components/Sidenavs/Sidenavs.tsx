import { usePathname } from 'next/navigation';
import { RootSidenav } from './RootSidenav';
import { OrganizationSidenav } from './OrganizationSidenav';

interface SidenavsProps {
  isOwner: boolean;
}

export const Sidenavs = (props: SidenavsProps) => {
  const pathname = usePathname();

  if (pathname?.includes('organizations')) {
    return <OrganizationSidenav />;
  }

  return <RootSidenav {...props} />;
};
