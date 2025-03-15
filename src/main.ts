import { Camera } from "./camera";
import { HittableList } from "./hittable-list";
import { Dialectric, Lambertian, Metal } from "./material";
import { Sphere } from "./sphere";
import { randomBetween, randomNum } from "./utils";
import { color, point3, Vec3 } from "./vec3";

function main() {
  const world = new HittableList();

  const ground = new Lambertian(color(0.5, 0.5, 0.5));
  world.add(new Sphere(point3(0, -1000, -1), 1000, ground));

  for (let a = -11; a < 11; a++) {
    for (let b = -11; b < 11; b++) {
      const chooseMaterial = randomNum();
      const center = point3(a + 0.9 * randomNum(), 0.2, b + 0.9 * randomNum());

      if (center.sub(point3(4, 0.2, 0)).length > 0.9) {
        if (chooseMaterial < 0.8) {
          const albedo = Vec3.random().vectorMultiply(Vec3.random());
          const material = new Lambertian(albedo);
          world.add(new Sphere(center, 0.2, material));
        } else if (chooseMaterial < 0.95) {
          const albedo = Vec3.randomBetween(0.5, 1);
          const fuzz = randomBetween(0, 0.5);
          const material = new Metal(albedo, fuzz);
          world.add(new Sphere(center, 0.2, material));
        } else {
          const material = new Dialectric(1.5);
          world.add(new Sphere(center, 0.2, material));
        }
      }
    }
  }
  const material1 = new Dialectric(1.5);
  world.add(new Sphere(point3(0, 1, 0), 1.0, material1));

  const material2 = new Lambertian(color(0.4, 0.2, 0.1));
  world.add(new Sphere(point3(-4, 1, 0), 1.0, material2));

  const material3 = new Metal(color(0.7, 0.6, 0.5), 0.0);
  world.add(new Sphere(point3(4, 1, 0), 1.0, material3));

  const cam = new Camera({
    width: 800,
    aspectRatio: 16 / 9,
    samplesPerPixel: 100,
    maxDepth: 50,
    vfov: 20,
    lookFrom: point3(13, 2, 3),
    lookAt: point3(0, 0, 0),
    defocusAngle: 0.6,
    focusDist: 10,
  });

  cam.render(world);
}

window.button.onclick = main;
