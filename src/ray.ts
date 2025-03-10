import { Vec3 } from "./vec3";

export function ray(origin: Vec3, direction: Vec3) {
  return new Ray(origin, direction);
}

export class Ray {
  constructor(
    public origin: Vec3,
    public dir: Vec3,
  ) {}

  at(t: number) {
    return this.origin.clone().add(this.dir.clone().k(t));
  }
}
