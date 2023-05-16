import * as React from 'react';
import { SVGProps } from 'react';
const SvgPinAltLight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    width={17}
    height={18}
    fill='none'
    stroke='#878787'
    {...props}
  >
    <path
      d='M13.714 12.542c.529.323.807.69.807 1.063s-.278.74-.807 1.063-1.288.591-2.204.778a15.34 15.34 0 0 1-3.01.285c-1.057 0-2.095-.098-3.01-.285s-1.675-.455-2.204-.778-.807-.689-.807-1.062.278-.739.807-1.062'
      strokeLinecap='round'
    />
    <path d='M13.813 7.584c0 3.434-3.629 5.788-4.888 6.502a.85.85 0 0 1-.849 0c-1.259-.714-4.888-3.068-4.888-6.502 0-3.187 2.574-5.312 5.313-5.312 2.833 0 5.313 2.125 5.313 5.313z' />
    <circle cx={8.499} cy={7.583} r={2.333} />
  </svg>
);
export default SvgPinAltLight;
