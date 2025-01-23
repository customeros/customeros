import { Play } from '@ui/media/icons/Play';
import { Target05 } from '@ui/media/icons/Target05';
import { Users02 } from '@ui/media/icons/Users02.tsx';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Shuffle01 } from '@ui/media/icons/Shuffle01';
import { Signature } from '@ui/media/icons/Signature';
import { Building07 } from '@ui/media/icons/Building07';
import { CheckHeart } from '@ui/media/icons/CheckHeart';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';

export const iconMap: Record<
  string,
  (props: React.SVGAttributes<SVGElement>) => JSX.Element
> = {
  InvoiceUpcoming: (props) => <InvoiceUpcoming {...props} />,
  InvoiceCheck: (props) => <InvoiceCheck {...props} />,
  Building07: (props) => <Building07 {...props} />,
  CheckHeart: (props) => <CheckHeart {...props} />,
  Users01: (props) => <Users02 {...props} />,
  Users02: (props) => <Users02 {...props} />,
  HeartHand: (props) => <HeartHand {...props} />,
  Signature: (props) => <Signature {...props} />,
  Target05: (props) => <Target05 {...props} />,
  CoinsStacked01: (props) => <CoinsStacked01 {...props} />,
  Shuffle01: (props) => <Shuffle01 {...props} />,
  Welcome: (props) => <Play {...props} />,
};
