export type RelationshipType =
  | 'CUSTOMER'
  | 'DISTRIBUTOR'
  | 'PARTNER'
  | 'LICENSING_PARTNER'
  | 'FRANCHISEE'
  | 'FRANCHISOR'
  | 'AFFILIATE'
  | 'RESELLER'
  | 'INFLUENCER_OR_CONTENT_CREATOR'
  | 'MEDIA_PARTNER'
  | 'INVESTOR'
  | 'MERGER_OR_ACQUISITION_TARGET'
  | 'PARENT_COMPANY'
  | 'SUBSIDIARY'
  | 'JOINT_VENTURE'
  | 'SPONSOR'
  | 'SUPPLIER'
  | 'VENDOR'
  | 'CONTRACT_MANUFACTURER'
  | 'ORIGINAL_EQUIPMENT_MANUFACTURER'
  | 'ORIGINAL_DESIGN_MANUFACTURER'
  | 'PRIVATE_LABEL_MANUFACTURER'
  | 'LOGISTICS_PARTNER'
  | 'CONSULTANT'
  | 'SERVICE_PROVIDER'
  | 'OUTSOURCING_PROVIDER'
  | 'INSOURCING_PARTNER'
  | 'TECHNOLOGY_PROVIDER'
  | 'DATA_PROVIDER'
  | 'CERTIFICATION_BODY'
  | 'STANDARDS_ORGANIZATION'
  | 'INDUSTRY_ANALYST'
  | 'REAL_ESTATE_PARTNER'
  | 'TALENT_ACQUISITION_PARTNER'
  | 'PROFESSIONAL_EMPLOYER_ORGANIZATION'
  | 'RESEARCH_COLLABORATOR'
  | 'REGULATORY_BODY'
  | 'TRADE_ASSOCIATION_MEMBER'
  | 'COMPETITOR';

export const relationshipOptions: { value: RelationshipType; label: string }[] =
  [
    {
      value: 'CUSTOMER',
      label: 'Customer',
    },
    {
      value: 'INVESTOR',
      label: 'Investor',
    },
    {
      value: 'VENDOR',
      label: 'Vendor',
    },
    {
      value: 'AFFILIATE',
      label: 'Affiliate',
    },
    {
      value: 'CERTIFICATION_BODY',
      label: 'Certification Body',
    },
    {
      value: 'COMPETITOR',
      label: 'Competitor',
    },
    {
      value: 'CONSULTANT',
      label: 'Consultant',
    },
    {
      value: 'CONTRACT_MANUFACTURER',
      label: 'Contract Manufacturer',
    },

    {
      value: 'DATA_PROVIDER',
      label: 'Data Provider',
    },
    {
      value: 'DISTRIBUTOR',
      label: 'Distributor',
    },
    {
      value: 'FRANCHISEE',
      label: 'Franchisee',
    },
    {
      value: 'FRANCHISOR',
      label: 'Franchisor',
    },
    {
      value: 'INDUSTRY_ANALYST',
      label: 'Industry Analyst',
    },
    {
      value: 'INFLUENCER_OR_CONTENT_CREATOR',
      label: 'Influencer or Content Creator',
    },
    {
      value: 'INSOURCING_PARTNER',
      label: 'Insourcing Partner',
    },
    {
      value: 'JOINT_VENTURE',
      label: 'Joint Venture',
    },
    {
      value: 'LICENSING_PARTNER',
      label: 'Licensing Partner',
    },
    {
      value: 'LOGISTICS_PARTNER',
      label: 'Logistics Partner',
    },
    {
      value: 'MEDIA_PARTNER',
      label: 'Media Partner',
    },
    {
      value: 'MERGER_OR_ACQUISITION_TARGET',
      label: 'Merger or Aquistion Target',
    },
    {
      value: 'ORIGINAL_DESIGN_MANUFACTURER',
      label: 'Original Design Manufacturer',
    },
    {
      value: 'ORIGINAL_EQUIPMENT_MANUFACTURER',
      label: 'Original Equipment Manufacturer',
    },
    {
      value: 'OUTSOURCING_PROVIDER',
      label: 'Outsourcing Provider',
    },
    {
      value: 'PARENT_COMPANY',
      label: 'Parent Company',
    },
    {
      value: 'PARTNER',
      label: 'Partner',
    },
    {
      value: 'PRIVATE_LABEL_MANUFACTURER',
      label: 'Private Label Manufacturer',
    },
    {
      value: 'PROFESSIONAL_EMPLOYER_ORGANIZATION',
      label: 'Professional Employer Organization',
    },
    {
      value: 'REAL_ESTATE_PARTNER',
      label: 'Real Estate Partner',
    },
    {
      value: 'REGULATORY_BODY',
      label: 'Regulatory Body',
    },
    {
      value: 'RESEARCH_COLLABORATOR',
      label: 'Research Collaborator',
    },
    {
      value: 'RESELLER',
      label: 'Reseller',
    },
    {
      value: 'SERVICE_PROVIDER',
      label: 'Service Provider',
    },
    {
      value: 'SPONSOR',
      label: 'Sponsor',
    },
    {
      value: 'STANDARDS_ORGANIZATION',
      label: 'Standards Organization',
    },
    {
      value: 'SUBSIDIARY',
      label: 'Subsidiary',
    },
    {
      value: 'SUPPLIER',
      label: 'Supplier',
    },
    {
      value: 'TALENT_ACQUISITION_PARTNER',
      label: 'Talent Aquisition Partner',
    },
    {
      value: 'TECHNOLOGY_PROVIDER',
      label: 'Technology Partner',
    },
    {
      value: 'TRADE_ASSOCIATION_MEMBER',
      label: 'Trade Association Member',
    },
  ];
