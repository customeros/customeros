import { ReactElement } from 'react';

import attio from '@assets/integrationOptionLogos/attio.svg';
import grain from '@assets/integrationOptionLogos/grain.webp';
import fathom from '@assets/integrationOptionLogos/fathom.webp';
import close from '@assets/integrationOptionLogos/close.com.svg';
import shopify from '@assets/integrationOptionLogos/shopify.webp';
import mixpanel from '@assets/integrationOptionLogos/mixpanel.svg';
import pipedrive from '@assets/integrationOptionLogos/pipedrive.webp';

import { DataSource } from '@graphql/types';
import { Stripe } from '@ui/media/logos/Stripe.tsx';
import { Outlook } from '@ui/media/logos/Outlook.tsx';
import { Zendesk } from '@ui/media/icons/Zendesk.tsx';
import { Intercom } from '@ui/media/icons/Intercom.tsx';
import { Salesforce } from '@ui/media/logos/Salesforce.tsx';

import logoCustomerOs from '../../../../../../../../assets/integrationOptionLogos/customer-os.png';

export const dataSourceLogo: Record<string, ReactElement> = {
  [DataSource.Openline]: (
    <img width={16} height={16} alt='CustomerOS' src={logoCustomerOs} />
  ),
  [DataSource.Grain]: <img width={16} height={16} alt='Grain' src={grain} />,
  [DataSource.Fathom]: <img width={16} height={16} alt='Fathom' src={fathom} />,
  [DataSource.Pipedrive]: (
    <img width={16} height={16} alt='Pipedrive' src={pipedrive} />
  ),
  [DataSource.Intercom]: <Intercom width={16} height={16} />,
  [DataSource.Close]: <img width={16} height={16} src={close} alt='Close' />,

  [DataSource.Mixpanel]: (
    <img width={16} height={16} src={mixpanel} alt='Mixpanel' />
  ),
  [DataSource.Shopify]: (
    <img width={16} height={16} src={shopify} alt='Shopify' />
  ),
  [DataSource.Salesforce]: <Salesforce width={16} height={16} />,
  [DataSource.Outlook]: <Outlook width={16} height={16} />,
  [DataSource.Stripe]: <Stripe width={16} height={16} />,
  [DataSource.ZendeskSell]: <Zendesk width={16} height={16} />,
  [DataSource.ZendeskSupport]: <Zendesk width={16} height={16} />,
  [DataSource.Attio]: <img width={16} height={16} src={attio} alt='Attio' />,
};
