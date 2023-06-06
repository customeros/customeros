import type { SVGProps } from 'react';

import Affiliate from '@spaces/atoms/icons/Affiliate';
import CertificateBody from '@spaces/atoms/icons/CertificationBody';
import Competitor from '@spaces/atoms/icons/Competitor';
import Consultant from '@spaces/atoms/icons/Consultant';
import ContractManufacturer from '@spaces/atoms/icons/ContractManufacturer';
import Customer from '@spaces/atoms/icons/Customer';
import DataProvider from '@spaces/atoms/icons/DataProvider';
import Distributor from '@spaces/atoms/icons/Distributor';
import Franchisee from '@spaces/atoms/icons/Franchisee';
import Franchisor from '@spaces/atoms/icons/Franchisor';
import IndustryAnalyst from '@spaces/atoms/icons/IndustryAnalyst';
import Influencer from '@spaces/atoms/icons/Influencer';
import InsourcingPartner from '@spaces/atoms/icons/InsourcingPartner';
import Investor from '@spaces/atoms/icons/Investor';
import JointVenture from '@spaces/atoms/icons/JointVenture';
import LicensingPartner from '@spaces/atoms/icons/LicensingPartner';
import LogisticPartner from '@spaces/atoms/icons/LogisticPartner';
import MediaPartner from '@spaces/atoms/icons/MediaPartner';
import MergerAcquisitionTarget from '@spaces/atoms/icons/MergerAcquisitionTarget';
import OriginalDesignManufacturer from '@spaces/atoms/icons/OriginalDesignManufacturer';
import OriginalEquipmentManufacturer from '@spaces/atoms/icons/OriginalEquipmentManufacturer';
import OutsourcingProvider from '@spaces/atoms/icons/OutsourcingProvider';
import ParentCompany from '@spaces/atoms/icons/ParentCompany';
import Partner from '@spaces/atoms/icons/Partner';
import PrivateLabelManufacturer from '@spaces/atoms/icons/PrivateLabelManufacturer';
import ProfessionalEmployerOrganization from '@spaces/atoms/icons/ProfessionalEmployerOrganization';
import RealEstatePartner from '@spaces/atoms/icons/RealEstatePartner';
import RegulatoryBody from '@spaces/atoms/icons/RegulatoryBody';
import ResearchBody from '@spaces/atoms/icons/ResearchBody';
import Reseller from '@spaces/atoms/icons/Reseller';
import ServiceProvider from '@spaces/atoms/icons/ServiceProvider';
import Sponsor from '@spaces/atoms/icons/Sponsor';
import StandardsOrganization from '@spaces/atoms/icons/StandardsOrganization';
import Subsidiary from '@spaces/atoms/icons/Subsidiary';
import Supplier from '@spaces/atoms/icons/Supplier';
import TalentAcquisitionPartner from '@spaces/atoms/icons/TalentAcquisitionPartner';
import ShieldTechnologyProvider from '@spaces/atoms/icons/ShieldTechnologyProvider';
import TradeAssociationMember from '@spaces/atoms/icons/TradeAssociationMember';
import Vendor from '@spaces/atoms/icons/Vendor';

import { RelationshipType } from './util';

export const icons: Record<
  RelationshipType,
  (props: SVGProps<SVGSVGElement>) => JSX.Element
> = {
  AFFILIATE: Affiliate,
  CERTIFICATION_BODY: CertificateBody,
  COMPETITOR: Competitor,
  CONSULTANT: Consultant,
  CONTRACT_MANUFACTURER: ContractManufacturer,
  CUSTOMER: Customer,
  DATA_PROVIDER: DataProvider,
  DISTRIBUTOR: Distributor,
  FRANCHISEE: Franchisee,
  FRANCHISOR: Franchisor,
  INDUSTRY_ANALYST: IndustryAnalyst,
  INFLUENCER_OR_CONTENT_CREATOR: Influencer,
  INSOURCING_PARTNER: InsourcingPartner,
  INVESTOR: Investor,
  JOINT_VENTURE: JointVenture,
  LICENSING_PARTNER: LicensingPartner,
  LOGISTICS_PARTNER: LogisticPartner,
  MEDIA_PARTNER: MediaPartner,
  MERGER_OR_ACQUISITION_TARGET: MergerAcquisitionTarget,
  ORIGINAL_DESIGN_MANUFACTURER: OriginalDesignManufacturer,
  ORIGINAL_EQUIPMENT_MANUFACTURER: OriginalEquipmentManufacturer,
  OUTSOURCING_PROVIDER: OutsourcingProvider,
  PARENT_COMPANY: ParentCompany,
  PARTNER: Partner,
  PRIVATE_LABEL_MANUFACTURER: PrivateLabelManufacturer,
  PROFESSIONAL_EMPLOYER_ORGANIZATION: ProfessionalEmployerOrganization,
  REAL_ESTATE_PARTNER: RealEstatePartner,
  REGULATORY_BODY: RegulatoryBody,
  RESEARCH_COLLABORATOR: ResearchBody,
  RESELLER: Reseller,
  SERVICE_PROVIDER: ServiceProvider,
  SPONSOR: Sponsor,
  STANDARDS_ORGANIZATION: StandardsOrganization,
  SUBSIDIARY: Subsidiary,
  SUPPLIER: Supplier,
  TALENT_ACQUISITION_PARTNER: TalentAcquisitionPartner,
  TECHNOLOGY_PROVIDER: ShieldTechnologyProvider,
  TRADE_ASSOCIATION_MEMBER: TradeAssociationMember,
  VENDOR: Vendor,
};

interface SelectMenuItemIconProps extends SVGProps<SVGSVGElement> {
  name: RelationshipType;
}

export const SelectMenuItemIcon = ({
  name,
  ...props
}: SelectMenuItemIconProps) => {
  return icons[name] ? icons[name](props) : null;
};
