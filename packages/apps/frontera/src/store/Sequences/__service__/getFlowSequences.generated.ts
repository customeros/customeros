import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type GetFlowSequencesQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetFlowSequencesQuery = {
  __typename?: 'Query';
  sequences: Array<{
    __typename?: 'FlowSequence';
    name: string;
    description: string;
    status: Types.FlowSequenceStatus;
    metadata: { __typename?: 'Metadata'; id: string };
    flow: {
      __typename?: 'Flow';
      name: string;
      description: string;
      status: Types.FlowStatus;
      metadata: { __typename?: 'Metadata'; id: string };
    };
    steps: Array<{
      __typename?: 'FlowSequenceStep';
      name: string;
      status: Types.FlowSequenceStepStatus;
      type: Types.FlowSequenceStepType;
      subtype?: Types.FlowSequenceStepSubtype | null;
      metadata: { __typename?: 'Metadata'; id: string };
    }>;
  }>;
};
