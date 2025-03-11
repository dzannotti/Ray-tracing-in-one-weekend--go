import { Ray } from "./ray";
import { Vec3 } from "./vec3";

export class HitRecord {
  public frontFace?: boolean;
  public p?: Vec3;
  public normal?: Vec3;
  public t?: number;

  setFaceNormal(r: Ray, outwardNormal: Vec3) {
    this.frontFace = Vec3.dot(r.direction, outwardNormal) < 0;
    this.normal = this.frontFace ? outwardNormal : outwardNormal.k(-1);
  }
}

export abstract class Hittable {
  abstract hit(
    r: Ray,
    rayTMin: number,
    rayTMax: number,
    rec: HitRecord,
  ): [boolean, HitRecord];
}
