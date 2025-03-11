export function interval(min: number, max: number) {
  return new Interval(min, max);
}

export class Interval {
  constructor(
    public min: number,
    public max: number,
  ) {}

  get size() {
    return this.max - this.min;
  }

  contains(x: number) {
    return this.min <= x && x <= this.max;
  }

  surrounds(x: number) {
    return this.min < x && x < this.max;
  }

  static empty = new Interval(Infinity, -Infinity);
  static universe = new Interval(-Infinity, +Infinity);
}
