'use client';
import { useRef, useEffect } from 'react';

import * as d3 from 'd3';

interface CircleData {
  x: number;
  y: number;
  radius: number;
}
interface UseBubblesOptions {
  width: number;
  height: number;
  data: CircleData[];
  transform?: string;
}

const TOP_R = 25; // max radius value
const LOW_R = 2; // min radius value

export const Bubbles = ({
  data,
  width,
  height,
  transform,
}: UseBubblesOptions) => {
  const tick = useRef(0);
  const topr = useRef(TOP_R);
  const ref = useRef<SVGGElement>(null);

  useEffect(() => {
    const simulation = makeSimulation({
      topr,
      tick,
      data,
      width,
      height,
      domRef: ref,
      rRange: [LOW_R, topr.current],
    });

    return () => {
      simulation.stop();
    };
  }, [data, width, height]);

  return <g ref={ref} width={width} height={height} transform={transform}></g>;
};

function getExtents(data: CircleData[]) {
  const minMaxX = d3.extent(data, (d) => d.x) as [number, number];
  const minMaxY = d3.extent(data, (d) => d.y) as [number, number];
  const minMaxR = d3.extent(data, (d) => d.radius) as [number, number];

  return {
    minMaxX,
    minMaxY,
    minMaxR,
  };
}

interface GetScalesOptions {
  width: number;
  height: number;
  rRange: [number, number];
  minMaxX: [number, number];
  minMaxY: [number, number];
  minMaxR: [number, number];
}
function getScales({
  width,
  height,
  rRange,
  minMaxX,
  minMaxY,
  minMaxR,
}: GetScalesOptions) {
  const xScale = d3.scaleLinear().domain(minMaxX).range([50, width]);
  const yScale = d3.scaleLinear().domain(minMaxY).range([20, height]);
  const rScale = d3.scaleSqrt().domain(minMaxR).range(rRange);

  return {
    xScale,
    yScale,
    rScale,
  };
}

function getScaledCoordinates(
  data: CircleData[],
  width: number,
  height: number,
  rRange: [number, number] = [LOW_R, TOP_R],
) {
  const { minMaxX, minMaxR, minMaxY } = getExtents(data);
  const { xScale, yScale, rScale } = getScales({
    width,
    height,
    minMaxX,
    minMaxY,
    minMaxR,
    rRange,
  });

  return {
    X: (d: CircleData) => xScale(d.x),
    Y: (d: CircleData) => yScale(d.y),
    R: (d: CircleData) => rScale(d.radius),
  };
}

function checkForCollisions(data: CircleData[]) {
  for (let i = 0; i < data.length; i++) {
    for (let j = i + 1; j < data.length; j++) {
      const circle1 = data[i];
      const circle2 = data[j];

      const dx = circle1.x - circle2.x;
      const dy = circle1.y - circle2.y;
      const distance = Math.sqrt(dx * dx + dy * dy);

      // Check if the circles overlap
      if (distance < circle1.radius + circle2.radius) {
        return true; // Collision detected
      }
    }
  }

  return false; // No collisions detected
}

interface MakeSimulationOptions {
  width: number;
  height: number;
  data: CircleData[];
  rRange: [number, number];
  topr: React.MutableRefObject<number>;
  tick: React.MutableRefObject<number>;
  domRef: React.MutableRefObject<SVGGElement | null>;
}

function makeSimulation(options: MakeSimulationOptions) {
  const { data, width, domRef, height, rRange, topr, tick } = options;

  const target = d3.select(domRef.current);
  const { X, Y, R } = getScaledCoordinates(data, width, height, rRange);

  const simulation = d3
    .forceSimulation<CircleData>(data)
    .force('x', d3.forceX<CircleData>((d) => X(d)).strength(0.1))
    .force('y', d3.forceY<CircleData>((d) => Y(d)).strength(0.1))
    .force('charge', d3.forceManyBody<CircleData>().strength(-20))
    .force(
      'collide',
      d3
        .forceCollide<CircleData>()
        .strength(1)
        .radius((d) => R(d))
        .iterations(10),
    );

  let out = simulation;

  target
    .selectAll<SVGCircleElement, CircleData>('circle')
    .data(data)
    .enter()
    .append('circle')
    .attr('cx', (d) => X(d))
    .attr('cy', (d) => Y(d))
    .attr('r', (d) => R(d))
    .style('fill', 'red')
    .style('opacity', 0.5);

  simulation.on('tick', () => {
    const _data = simulation.nodes();
    const { X, Y, R } = getScaledCoordinates(_data, width, height, rRange);

    const hasCollisions = checkForCollisions(
      simulation.nodes().map((n) => ({ ...n, x: X(n), y: Y(n), radius: R(n) })),
    );

    d3.select(domRef.current)
      .selectAll('circle')
      .data(data)
      .join('circle')
      .attr('cx', (d) => X(d))
      .attr('cy', (d) => Y(d))
      .attr('r', (d) => R(d));

    if (tick.current > 3 && hasCollisions) {
      simulation.stop();

      topr.current -= 1;
      tick.current = 0;

      out = makeSimulation({
        ...options,
        rRange: [LOW_R, topr.current],
      });
    }
    tick.current += 1;
  });

  return out;
}
