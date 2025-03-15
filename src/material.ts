import { HitRecord } from "./hittable";
import { ray, Ray } from "./ray";
import { Ref } from "./ref";
import { randomNum } from "./utils";
import { color, Vec3 } from "./vec3";

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

export class Dialectric extends Material {
  constructor(public refractionIndex: number) {
    super();
  }

  scatter(
    rayIn: Ref<Ray>,
    rec: HitRecord,
    attenuation: Ref<Vec3>,
    scattered: Ref<Ray>,
  ): boolean {
    attenuation.value = color(1, 1, 1);

    const ri = rec.frontFace!
      ? 1.0 / this.refractionIndex
      : this.refractionIndex;

    const unitDirection = rayIn.value!.direction.unit;

    const cosTheta = Math.min(Vec3.dot(unitDirection.k(-1), rec.normal!), 1.0);
    const sinTheta = Math.sqrt(1 - cosTheta * cosTheta);

    const cannotRefract = ri * sinTheta > 1;

    let direction: Vec3;

    if (cannotRefract || this.reflectance(cosTheta, ri) > randomNum()) {
      direction = Vec3.reflect(unitDirection, rec.normal!);
    } else {
      direction = Vec3.refract(unitDirection, rec.normal!, ri);
    }

    scattered.value = ray(rec.p!, direction);

    return true;
  }

  reflectance(cosine: number, refractionIndex: number) {
    // Schlick's approximation
    let r0 = (1 - refractionIndex) / (1 + refractionIndex);
    r0 = r0 * r0;

    return r0 + (1 - r0) * Math.pow(1 - cosine, 5);
  }
}
