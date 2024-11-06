import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type SaveCustomFieldTemplateMutationVariables = Types.Exact<{
  input: Types.CustomFieldTemplateInput;
}>;

export type SaveCustomFieldTemplateMutation = {
  __typename?: 'Mutation';
  customFieldTemplate_Save: {
    __typename?: 'CustomFieldTemplate';
    id: string;
    name: string;
    type: Types.CustomFieldTemplateType;
    createdAt: any;
    updatedAt: any;
  };
};
