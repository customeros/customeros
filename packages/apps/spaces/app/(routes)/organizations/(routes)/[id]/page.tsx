import { SideSection } from './components/SideSection';
import { MainSection } from './components/MainSection';
import {
  OrganizationInfo,
  OrganizationLogo,
  OrganizationTabs,
  OrganizationHeader,
} from './components/OrganizationInfo';

export default function OrganizationPage() {
  return (
    <>
      <SideSection>
        <OrganizationInfo>
          <OrganizationHeader>
            <OrganizationLogo src='/logos/bigquery.svg' />
          </OrganizationHeader>

          <OrganizationTabs />
        </OrganizationInfo>
      </SideSection>
      <MainSection></MainSection>
    </>
  );
}
