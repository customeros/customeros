export function dodgeBottom(data, { radius = 1, x = (d) => d } = {}) {
  const radius2 = radius ** 2;
  const circles = data
    .map((d, i, data) => ({ x: +x(d, i, data), data: d }))
    .sort((a, b) => a.x - b.x);
  const epsilon = 1e-3;
  let head = null,
    tail = null;

  // Returns true if circle ⟨x,y⟩ intersects with any circle in the queue.
  function intersects(x, y) {
    let a = head;
    while (a) {
      if (radius2 - epsilon > (a.x - x) ** 2 + (a.y - y) ** 2) {
        return true;
      }
      a = a.next;
    }

    return false;
  }

  // Place each circle sequentially.
  for (const b of circles) {
    // Remove circles from the queue that can’t intersect the new circle b.
    while (head && head.x < b.x - radius2) head = head.next;

    // Choose the minimum non-intersecting tangent.
    if (intersects(b.x, (b.y = 0))) {
      let a = head;
      b.y = Infinity;
      do {
        let y = a.y + Math.sqrt(radius2 - (a.x - b.x) ** 2);
        if (y < b.y && !intersects(b.x, y)) b.y = y;
        a = a.next;
      } while (a);
    }

    // Add b to the queue.
    b.next = null;
    if (head === null) head = tail = b;
    else tail = tail.next = b;
  }

  return circles;
}

export function dodgeMiddle(data, { radius, x }) {
  // const radius2 = radius ** 2;
  const circles = data
    .map((d) => ({ x: x(d), data: d }))
    .sort((a, b) => a.x - b.x);
  const epsilon = 1e-3;
  let head = null,
    tail = null;

  // Returns true if circle ⟨x,y⟩ intersects with any circle in the queue.
  function intersects(x, y) {
    let a = head;
    while (a) {
      const r = Math.max(radius(a.data), radius({ x, data: null }));
      const r2 = r ** 2;
      if (r2 - epsilon > (a.x - x) ** 2 + (a.y - y) ** 2) {
        return true;
      }
      a = a.next;
    }

    return false;
  }

  // Place each circle sequentially.
  for (const b of circles) {
    const radius2 = radius(b.data) ** 2;
    // Remove circles from the queue that can’t intersect the new circle b.
    while (head && head.x < b.x - radius2) head = head.next;

    // Choose the minimum non-intersecting tangent.
    if (intersects(b.x, (b.y = 0))) {
      let a = head;
      b.y = Infinity;
      do {
        let y1 = a.y + Math.sqrt(radius2 - (a.x - b.x) ** 2);
        let y2 = a.y - Math.sqrt(radius2 - (a.x - b.x) ** 2);
        if (Math.abs(y1) < Math.abs(b.y) && !intersects(b.x, y1)) b.y = y1;
        if (Math.abs(y2) < Math.abs(b.y) && !intersects(b.x, y2)) b.y = y2;
        a = a.next;
      } while (a);
    }

    // Add b to the queue.
    b.next = null;
    if (head === null) head = tail = b;
    else tail = tail.next = b;
  }

  return circles;
}

export function dodgeVariable(data, { x, r, padding } = {}) {
  const circles = data
    .map((d) => ({ x: x(d), r: r(d), data: d }))
    .sort((a, b) => b.r - a.r);
  const epsilon = 1e-3;
  let head = null,
    tail = null,
    _queue = null;

  // Returns true if circle ⟨x,y⟩ intersects with any circle in the queue.
  function intersects(x, y, r) {
    let a = head;
    while (a) {
      const radius2 = (a.r + r + padding) ** 2;
      if (radius2 - epsilon > (a.x - x) ** 2 + (a.y - y) ** 2) {
        return true;
      }
      a = a.next;
    }

    return false;
  }

  // Place each circle sequentially.
  for (const b of circles) {
    // Choose the minimum non-intersecting tangent.
    if (intersects(b.x, (b.y = b.r), b.r)) {
      let a = head;
      b.y = Infinity;
      do {
        let y = a.y + Math.sqrt((a.r + b.r + padding) ** 2 - (a.x - b.x) ** 2);
        if (y < b.y && !intersects(b.x, y, b.r)) b.y = y;
        a = a.next;
      } while (a);
    }

    // Add b to the queue.
    b.next = null;
    if (head === null) {
      head = tail = b;
      _queue = head;
    } else tail = tail.next = b;
  }

  return circles;
}
