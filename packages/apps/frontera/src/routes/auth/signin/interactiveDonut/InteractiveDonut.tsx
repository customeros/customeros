import React, { useRef, useEffect } from 'react';

import './style.scss';

export const InteractiveDonut = () => {
  const pathRef = useRef<SVGPathElement>(null);
  const gradientDivRef = useRef<HTMLDivElement>(null);
  const cosRingRef = useRef<SVGSVGElement>(null);

  const composeCircle = (
    radius: number,
    centerX: number,
    centerY: number,
    largeArcFlag: number,
  ): string => {
    return `M ${centerX} ${
      centerY - radius
    } A ${radius} ${radius} 0 1 ${largeArcFlag} ${centerX} ${
      centerY + radius
    } A ${radius} ${radius} 0 1 ${largeArcFlag} ${centerX} ${
      centerY - radius
    } Z`;
  };

  useEffect(() => {
    const updateDonutPath = (event?: MouseEvent) => {
      if (!pathRef.current || !gradientDivRef.current || !cosRingRef.current)
        return;

      const centerX = window.innerWidth / 2;
      const centerY = window.innerHeight / 2;
      const mouseX = event ? event.clientX : centerX;
      const mouseY = event ? event.clientY : centerY;
      const xOffset = mouseX - centerX;
      const yOffset = mouseY - centerY;
      // const xRatio = xOffset / centerX;
      // const yRatio = yOffset / centerY;
      const innerCircleX = 125 + xOffset * 0.025;
      const innerCircleY = 128 + yOffset * 0.015;
      const outerCircleX = 125 + xOffset * -0.025;
      const outerCircleY = 120 + yOffset * -0.015;

      const numDonuts = 16;
      const thickness = 8;
      const spacing = 2;

      let donutsPath = '';

      for (let i = 0; i < numDonuts; i++) {
        const outerRadius = 50 + Math.pow(i * (thickness + spacing), 0.9);
        const innerRadius = outerRadius - thickness + Math.pow(i, 0.8);
        const largeArcFlag = i % 2;
        const circleX =
          innerCircleX + (outerCircleX - innerCircleX) * (i / numDonuts);
        const circleY =
          innerCircleY + (outerCircleY - innerCircleY) * (i / numDonuts);

        donutsPath +=
          composeCircle(outerRadius, circleX, circleY, largeArcFlag) + ' ';
        donutsPath +=
          composeCircle(innerRadius, circleX, circleY, 1 - largeArcFlag) + ' ';
      }

      pathRef.current.setAttribute('d', donutsPath.trim());

      const angle = (xOffset / centerX) * 30;
      const yPos = 48 - (yOffset / centerY) * 5;

      gradientDivRef.current.style.background = `conic-gradient(from ${
        180 + angle
      }deg at 50% ${yPos}%, #4C375A 1deg, #464068 2deg, #294868 6deg, #185070 12deg, #1F688C 32deg, #21759A 40deg, #3DA5BE 90deg, #ACD0D9 135deg, #EFF6F8 154deg, #FFF 160deg, #FFF4D7 163deg, #FEE7A6 170deg, #F7BE33 178deg, #DE7324 182deg, #E9571F 220deg, #DF353C 260deg, #BF1B4E 335deg, #8E2C56 356deg, #4C375A 360deg)`;

      const shadowX = (mouseX - centerX) * -0.05;
      const shadowY = (mouseY - centerY) * -0.05;

      cosRingRef.current.style.filter = `drop-shadow(${shadowX}rem ${shadowY}rem 0.8rem #00000022)`;
    };

    document.addEventListener('mousemove', updateDonutPath);
    updateDonutPath();

    return () => {
      document.removeEventListener('mousemove', updateDonutPath);
    };
  }, []);

  return (
    <div
      style={{ perspective: '1000px' }}
      className='h-[300vh]  p-0 m-0 overflow-hidden donut relative z-999'
    >
      <svg
        fill='none'
        id='cos-ring'
        ref={cosRingRef}
        viewBox='0 0 265 260'
        className='w-full h-full'
        xmlns='http://www.w3.org/2000/svg'
        style={{
          transform: 'rotate3d(1, -1, 0, 20deg)',
          animation:
            'wobble 22s cubic-bezier(0.68, -0.55, 0.265, 1.55) infinite',
          transformStyle: 'preserve-3d',
          overflow: 'visible',
          filter: 'drop-shadow(0 20rem 1.5rem #00000033)',
        }}
      >
        <foreignObject
          x='-30'
          y='-30'
          width='320'
          height='320'
          clipPath='url(#cos-donut)'
          style={{ overflow: 'visible', transform: 'translateXYZ(14%,14%,0)' }}
        >
          <div
            ref={gradientDivRef}
            className='w-full h-full overflow-visible'
            style={{
              background:
                'conic-gradient(from 180deg at 50% 64%, #4c375a 1deg, #464068 2deg, #294868 6deg, #185070 12deg, #1f688c 32deg, #21759a 40deg, #3da5be 90deg, #acd0d9 135deg, #eff6f8 154deg, #fff 160deg, #fff4d7 163deg, #fee7a6 170deg, #f7be33 178deg, #de7324 182deg, #e9571f 220deg, #df353c 260deg, #bf1b4e 335deg, #8e2c56 356deg, #4c375a 360deg)',
              transition: 'background 0.5s ease-out',
            }}
          />
        </foreignObject>
        <defs>
          <clipPath id='cos-donut'>
            <path ref={pathRef} />
          </clipPath>
        </defs>
      </svg>
    </div>
  );
};
