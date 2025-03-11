import { HitRecord, Hittable } from "./hittable";
import { Ray } from "./ray";
import { Vec3 } from "./vec3";

export class Sphere extends Hittable {
  constructor(
    public center: Vec3,
    public radius: number,
  ) {
    super();
  }

  hit(
    r: Ray,
    rayTMin: number,
    rayTMax: number,
    rec: HitRecord,
  ): [boolean, HitRecord] {
    const oc = this.center.sub(r.origin);
    const a = r.direction.lengthSquared;
    const h = Vec3.dot(r.direction, oc);
    const c = oc.lengthSquared - this.radius * this.radius;
    const discriminant = h * h - a * c;

    if (discriminant < 0) {
      return [false, rec];
    }

    let root = (h - Math.sqrt(discriminant)) / a;

    if (root <= rayTMin || rayTMax <= root) {
      root = (h + Math.sqrt(discriminant)) / a;

      if (root <= rayTMin || rayTMax <= root) {
        return [false, rec] as const;
      }
    }

    rec.t = root;
    rec.p = r.at(rec.t);
    let outwardNormal = rec.p.sub(this.center).div(this.radius);
    rec.setFaceNormal(r, outwardNormal);

    return [true, rec] as const;
  }
}
