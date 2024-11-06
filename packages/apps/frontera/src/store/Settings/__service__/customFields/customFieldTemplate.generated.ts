import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetCustomFieldTemplatesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetCustomFieldTemplatesQuery = {
  __typename?: 'Query';
  customFieldTemplate_List: Array<{
    __typename?: 'CustomFieldTemplate';
    id: string;
    name: string;
    type: Types.CustomFieldTemplateType;
    createdAt: any;
    updatedAt: any;
    entityType: Types.EntityType;
    validValues: Array<string>;
    order?: any | null;
    required?: boolean | null;
    length?: any | null;
    min?: any | null;
    max?: any | null;
  }>;
};
