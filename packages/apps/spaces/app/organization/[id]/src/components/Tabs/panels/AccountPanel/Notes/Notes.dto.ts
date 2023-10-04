import { OrganizationAccountDetailsQuery } from '@organization/src/graphql/getAccountPanelDetails.generated';
import { UpdateOrganizationMutationVariables } from '@organization/src/graphql/updateOrganization.generated';

export interface NotesForm {
  id: string;
  note: string;
}

export class NotesDTO implements NotesForm {
  id: string;
  note: string;

  constructor(
    data?: Partial<OrganizationAccountDetailsQuery['organization']> | null,
  ) {
    this.id = data?.id || '';
    this.note = data?.note || '';
  }

  static toForm(data: OrganizationAccountDetailsQuery) {
    return new NotesDTO(data.organization);
  }

  static toPayload(data: NotesForm) {
    return {
      id: data.id,
      note: data.note,
      name: '',
      patch: true,
    } as UpdateOrganizationMutationVariables['input'];
  }
}
