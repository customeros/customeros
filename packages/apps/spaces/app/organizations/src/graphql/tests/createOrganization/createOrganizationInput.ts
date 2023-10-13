import { CustomFieldDataType, Market, OrganizationInput } from '../../../../../../app/src/types/__generated__/graphql.types';

export const createOrganizationInputVariables: OrganizationInput = {
  // referenceId: 'createOrganizationNameReferenceId',
  name: 'createOrganizationName',
  // description: 'createOrganizationDescription',
  // note: 'createOrganizationNote',
  domains: ['createOrganizationDomain_1', 'createOrganizationDomain_2'],
  // website: 'www.createOrganizationWebsite.com',
  // industry: 'createOrganizationIndustry',
  // subIndustry: 'createOrganizationSubindustry',
  // industryGroup: 'createOrganizationIndustryGroup',
  // isPublic: true,
  // isCustomer: true,
  customFields: [
    {
      name: 'customFields_1',
      datatype: CustomFieldDataType.Text,
      value: 'customFieldsValue_1',
    },
    {
      name: 'customFields_2',
      datatype: CustomFieldDataType.Bool,
      value: true,
    },
  ],
  fieldSets: [
    {
      name: 'fieldSets_1',
      customFields: [
        {
          name: 'customFields_1',
          datatype: CustomFieldDataType.Text,
          value: 'customFieldsValue_1',
        },
        {
          name: 'customFields_2',
          datatype: CustomFieldDataType.Bool,
          value: true,
        },
      ],
    },
    {
      name: 'fieldSets_2',
      customFields: [
        {
          name: 'customFields_2',
          datatype: CustomFieldDataType.Bool,
          value: true,
        },
      ],
    },
  ],
  // templateId: 'createOrganizationTemplateId',
  // market: Market.B2B,
  // employees: 2,
  // appSource: 'createOrganizationAppsource',
}