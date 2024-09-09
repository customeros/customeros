import React from 'react';

import { Target05 } from '@ui/media/icons/Target05';
import { Users01 } from '@ui/media/icons/Users01.tsx';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Building07 } from '@ui/media/icons/Building07';
import { CheckHeart } from '@ui/media/icons/CheckHeart';
import { Shuffle01 } from '@ui/media/icons/Shuffle01.tsx';
import { Signature } from '@ui/media/icons/Signature.tsx';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';

export const iconMap: Record<
  string,
  (props: React.SVGAttributes<SVGElement>) => JSX.Element
> = {
  InvoiceUpcoming: (props) => <InvoiceUpcoming {...props} />,
  InvoiceCheck: (props) => <InvoiceCheck {...props} />,
  Building07: (props) => <Building07 {...props} />,
  CheckHeart: (props) => <CheckHeart {...props} />,
  Users01: (props) => <Users01 {...props} />,
  HeartHand: (props) => <HeartHand {...props} />,
  Signature: (props) => <Signature {...props} />,
  Target05: (props) => <Target05 {...props} />,
  CoinsStacked01: (props) => <CoinsStacked01 {...props} />,
  Shuffle01: (props) => <Shuffle01 {...props} />,
};
