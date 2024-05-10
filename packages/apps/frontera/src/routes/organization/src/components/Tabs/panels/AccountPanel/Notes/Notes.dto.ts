import { UpdateOrganizationMutationVariables } from '@shared/graphql/updateOrganization.generated';
import { OrganizationAccountDetailsQuery } from '@organization/graphql/getAccountPanelDetails.generated';

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
      patch: true,
    } as UpdateOrganizationMutationVariables['input'];
  }
}
