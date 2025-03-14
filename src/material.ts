import { HitRecord } from "./hittable";
import { ray, Ray } from "./ray";
import { Ref } from "./ref";
import { Vec3 } from "./vec3";

export abstract class Material {
  abstract scatter(
    rayIn: Ref<Ray>,
    rec: HitRecord,
    attenuation: Ref<Vec3>,
    scattered: Ref<Ray>,
  ): boolean;
}

export class Lambertian extends Material {
  constructor(public albedo: Vec3) {
    super();
  }

  override scatter(
    _rayIn: Ref<Ray>,
    rec: HitRecord,
    attenuation: Ref<Vec3>,
    scattered: Ref<Ray>,
  ): boolean {
    let scatterDirection = rec.normal!.add(Vec3.randomUnitVector());

    if (scatterDirection.nearZero()) {
      scatterDirection = rec.normal!.clone();
    }

    scattered.value = ray(rec.p!, scatterDirection);
    attenuation.value = this.albedo;
    return true;
  }
}

export class Metal extends Material {
  fuzz: number;
  constructor(
    public albedo: Vec3,
    fuzz: number,
  ) {
    super();

    this.fuzz = fuzz < 1 ? fuzz : 1;
  }

  override scatter(
    rayIn: Ref<Ray>,
    rec: HitRecord,
    attenuation: Ref<Vec3>,
    scattered: Ref<Ray>,
  ): boolean {
    let reflected = Vec3.reflect(rayIn.value!.direction, rec.normal!);
    reflected = reflected.unit.add(Vec3.randomUnitVector().k(this.fuzz));
    scattered.value = ray(rec.p!, reflected);
    attenuation.value = this.albedo;
    return Vec3.dot(scattered.value.direction!, rec.normal!) > 0;
  }
}
