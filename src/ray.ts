import { Vec3 } from "./vec3";

export function ray(origin: Vec3, direction: Vec3) {
  return new Ray(origin, direction);
}

export class Ray {
  _origin: Vec3;
  _dir: Vec3;

  constructor(origin: Vec3, dir: Vec3) {
    this._origin = origin;
    this._dir = dir;
  }

  get origin() {
    return this._origin.clone();
  }

  get direction() {
    return this._dir.clone();
  }

  at(t: number) {
    return this.origin.add(this.direction.k(t));
  }
}
