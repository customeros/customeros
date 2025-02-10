import { FlowStatuses } from './pages/flows/flowsStatuses';

export const companies = {
  create: {
    name: '',
    domain: 'gomag.ro',
    orgAboutDescription: 'This is an initial description',
    orgAboutTags: 'testOrgTag',
    orgAboutRelationship: 'Not a fit',
    orgAboutRelationshipRequest: 'NOT_A_FIT',
    orgAboutIndustry: 'Independent Power and Renewable Electricity Producers',
    orgAboutBusinessType: 'B2C',
    orgAboutLastFundingRound: 'FRIENDS_AND_FAMILY',
    orgAboutNumberOfEmployees: '10000+ employees',
    orgAboutNumberOfEmployeesRequest: 10001,
    orgAboutOwner: 'customeros.fe.testing',
    orgAboutSocialLinkEmpty: '/qweasdzxc123ads',
  },
  update: {
    name: 'Yahoo! Inc.',
    domain: 'protv.ro',
    orgAboutDescription:
      'This company is simply the best, better than all the rest',
    orgAboutTags: 'testOrgTag',
    orgAboutRelationship: 'Not a fit',
    orgAboutRelationshipRequest: 'NOT_A_FIT',
    orgAboutIndustry: 'Independent Power and Renewable Electricity Producers',
    orgAboutBusinessType: 'B2C',
    orgAboutLastFundingRound: 'Friends and Family',
    orgAboutLastFundingRoundRequest: 'FRIENDS_AND_FAMILY',
    orgAboutNumberOfEmployees: '10000+ employees',
    orgAboutNumberOfEmployeesRequest: 10001,
    orgAboutOwner: 'customeros.fe.testing',
    orgAboutSocialLinkEmpty: '/qweasdzxc123ads',
  },
};

export const flow = {
  create: {
    status: FlowStatuses.Stopped,
    onHold: '0',
    ready: '0',
    scheduled: '0',
    inProgress: '0',
    completed: '0',
    goalAchieved: '0',
  },
  update: {
    status: FlowStatuses.Stopped,
    onHold: '0',
    ready: '1',
    scheduled: '0',
    inProgress: '0',
    completed: '0',
    goalAchieved: '0',
  },
};
